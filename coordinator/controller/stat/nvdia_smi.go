package stat

import (
	"log"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func NvidiaSMI() (string, error) {
	cmd := exec.Command("nvidia-smi")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func NvidiaSMIHandler(c *gin.Context) {
	output, err := NvidiaSMI()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Failed to connect to nvidia daemons"})
		return
	}
	c.String(200, output)
}
