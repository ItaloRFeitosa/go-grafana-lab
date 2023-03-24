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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/italorfeitosa/go-grafana-lab/chaos"
	"github.com/italorfeitosa/go-grafana-lab/checkout"
	"github.com/italorfeitosa/go-grafana-lab/order"
	"github.com/italorfeitosa/go-grafana-lab/payment"
	"github.com/italorfeitosa/go-grafana-lab/tracing"
	"github.com/italorfeitosa/go-grafana-lab/warehouse"

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
	checkoutClient := checkout.NewClient()
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
			"isSarehouse": isWarehouse,
		})
	})

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
	})
}

func GracefulShutdown(shutdownCallback func(context.Context)) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	log.Println("gracefully shutdown process...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer signal.Stop(quit)

	shutdownCallback(ctx)

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}

	log.Println("process exiting")
}
