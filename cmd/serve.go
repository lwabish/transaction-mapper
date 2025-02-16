package cmd

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lwabish/transaction-mapper/pkg/bank"
	"github.com/lwabish/transaction-mapper/pkg/consumer"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"path"
	"time"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()

		config := cors.DefaultConfig()
		config.AllowAllOrigins = true
		config.AddAllowHeaders("Content-Disposition")

		r.Use(gin.Recovery(), cors.New(config))
		api := r.Group("/api/v1")
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"Hello": "transaction mapper server"})
		})
		api.GET("/banks", func(c *gin.Context) {
			c.JSON(200, gin.H{"data": bank.Registry.List()})
		})
		api.GET("/apps", func(c *gin.Context) {
			c.JSON(200, gin.H{"data": consumer.Registry.List()})
		})
		api.POST("/transform", func(c *gin.Context) {
			formFile, err := c.FormFile("input")
			if err != nil {
				log.Printf("get form file error: %s", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			f, err := formFile.Open()
			if err != nil {
				log.Printf("open form file error: %s", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer func() {
				if err := f.Close(); err != nil {
					log.Printf("close form file error: %s", err)
				}
			}()
			content, err := io.ReadAll(f)
			if err != nil {
				log.Printf("read form file error: %s", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			arg := defaultArg()
			if err := c.ShouldBind(&arg); err != nil {
				log.Printf("bind form data error: %s", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			dstFile := fmt.Sprintf("%s-%d-%s.csv", path.Base(formFile.Filename), time.Now().Unix(), arg.App)
			if err := runEngine(arg, content, dstFile); err != nil {
				log.Printf("run engine error: %s", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.Header("Access-Control-Expose-Headers", "Content-Disposition")
			c.FileAttachment(dstFile, dstFile)
		})
		if err := r.Run(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
