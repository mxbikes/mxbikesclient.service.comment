package main

import (
	"net"
	"os"

	"github.com/mxbikes/mxbikesclient.service.comment/handler"
	"github.com/mxbikes/mxbikesclient.service.comment/repository"
	protobuffer "github.com/mxbikes/protobuf/comment"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	logLevel    = GetEnv("LOG_LEVEL", "info")
	port        = GetEnv("PORT", ":4089")
	postgresUrl = GetEnv("POSTGRES_URI", "host=host.docker.internal port=5432 user=postgres password=password sslmode=disable timezone=UTC connect_timeout=5")
)

func main() {
	logger := &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &prefixed.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}

	/* Database */
	db, err := gorm.Open(postgres.Open(postgresUrl), &gorm.Config{})
	if err != nil {
		logger.WithFields(logrus.Fields{"prefix": "POSTGRES"}).Fatal("unable to open a connection to database")
	}
	logger.WithFields(logrus.Fields{"prefix": "POSTGRES"}).Info("connection has been established successfully!")
	repo := repository.NewRepository(db)
	repo.Migrate()

	/* Server */
	// Create a tcp listener
	listener, err := net.Listen("tcp", port)
	if err != nil {
		logger.WithFields(logrus.Fields{"prefix": "SERVICE.MOD"}).Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	protobuffer.RegisterCommentServiceServer(grpcServer, handler.New(repo, *logger))
	reflection.Register(grpcServer)

	// Start grpc server on listener
	logger.WithFields(logrus.Fields{"prefix": "SERVICE.MOD"}).Infof("is listening on Grpc PORT: {%v}", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		logger.WithFields(logrus.Fields{"prefix": "SERVICE.MOD"}).Fatalf("failed to serve: %v", err)
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
