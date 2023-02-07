package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tanya-lyubimaya/WB_L0/internal/domain"
	"net/http"
	"strings"
	"time"
)

var (
	origins = map[string]bool{
		"http://localhost": true,
		"http://127.0.0.1": true,
	}
)

type Server struct {
	srv    *http.Server
	uc     domain.UseCase
	logger *logrus.Entry
}

func New(uc domain.UseCase, logger *logrus.Entry) (*Server, error) {
	router := gin.Default()
	router.Use(CORSMiddleware())
	server := &Server{uc: uc, srv: &http.Server{Handler: router}, logger: logrus.NewEntry(logger.Logger)}

	router.GET("/orders/:id", server.GetOrderByID)
	gin.LoggerWithWriter(logger.Writer())
	return server, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Request.Referer()) != 0 {
			for v := range origins {
				if strings.Index(v, c.Request.Referer()) > -1 {
					if c.Request.Referer()[len(c.Request.Referer())-1:] == "/" {
						c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Referer()[0:len(c.Request.Referer())-1])
					} else {
						c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Referer())
					}
					break
				}
			}
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Referer")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) Serve(port string) error {
	s.srv.Addr = port
	return s.srv.ListenAndServe()
}

func (s *Server) GracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	go func() {
		err := s.srv.Shutdown(ctx)
		if err != nil {
			s.logger.Fatalln(err)
		}
	}()
	<-ctx.Done()
	s.uc.Close()
}

func (s *Server) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		s.logger.Errorln("Invalid ID")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	result, err := s.uc.ReadOrderByID(id)
	if err != nil {
		s.logger.Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}
