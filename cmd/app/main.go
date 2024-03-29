package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/italorfeitosa/go-grafana-lab/internal/checkout"
	checkoutclient "github.com/italorfeitosa/go-grafana-lab/internal/checkout/client"
	"github.com/italorfeitosa/go-grafana-lab/internal/order"
	"github.com/italorfeitosa/go-grafana-lab/internal/payment"
	"github.com/italorfeitosa/go-grafana-lab/internal/warehouse"
	"github.com/italorfeitosa/go-grafana-lab/pkg/chaos"
	"github.com/italorfeitosa/go-grafana-lab/pkg/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/google/uuid"

	"github.com/gofiber/contrib/otelfiber"
)

var (
	isCheckout  bool
	isPayment   bool
	isOrder     bool
	isWarehouse bool
	isChaos     bool
)

func init() {
	flag.BoolVar(&isCheckout, "checkout", false, "run checkout")
	flag.BoolVar(&isPayment, "payment", false, "run payment")
	flag.BoolVar(&isOrder, "order", false, "run order")
	flag.BoolVar(&isWarehouse, "warehouse", false, "run warehouse")
	flag.BoolVar(&isChaos, "chaos", false, "run chaos")
}

func main() {
	flag.Parse()

	if isChaos {
		startChaos()
		return
	}

	startServer()
}

func startChaos() {
	checkoutClient := checkoutclient.New()
	stop := chaos.Do(func() {
		log.Println("calling checkout")
		id := uuid.NewString()
		if err := checkoutClient.StartCheckout(context.Background(), id); err != nil {
			log.Printf("error: %+v", err)
			return
		}

		err := checkoutClient.StartCheckoutPayment(context.Background(), id)
		log.Printf("error: %+v", err)

	}, 10*time.Second)

	GracefulShutdown(func(ctx context.Context) {
		stop()
	})
}

func startServer() {
	tp := tracing.SetupJaegerProvider()
	defer tp.Shutdown(context.Background())

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println(err.Error())
			return fiber.DefaultErrorHandler(c, err)
		},
	})

	app.Use(recover.New())

	app.Use(logger.New())

	app.Use(otelfiber.Middleware())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"isCheckout":  isCheckout,
			"isPayment":   isPayment,
			"isOrder":     isOrder,
			"isWarehouse": isWarehouse,
		})
	})

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	if isCheckout {
		checkout.SetRoutes(app)
	}

	if isWarehouse {
		warehouse.SetRoutes(app)
	}

	if isPayment {
		payment.SetRoutes(app)
	}

	if isOrder {
		order.SetRoutes(app)
	}

	err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))

	if err != nil {
		log.Fatal(err)
	}

	GracefulShutdown(func(ctx context.Context) {
		if err := app.Shutdown(); err != nil {
			log.Fatal(err)
		}
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	})
}

func GracefulShutdown(shutdownCallback func(context.Context)) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	log.Println("gracefully shutdown process...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer signal.Stop(quit)

	go shutdownCallback(ctx)

	<-ctx.Done()

	log.Println("process exiting")
}
