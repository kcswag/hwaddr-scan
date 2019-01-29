package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

type Result struct{
	IP string
	hostname string
	mac string
}

var (
	resSlice []Result
	resSliceIndex int
	ipRange string
)

var rootCmd = &cobra.Command{Use: "hwaddr",}

var scanCmd = &cobra.Command{
	Use: "scan",
	Short:"Scan the hardware address of devices in the same LAN",
	DisableFlagParsing:true,
	Run: func(cmd *cobra.Command, args []string) {
		//check if hwaddr-scan has been installed
		_,checkErr := runCommand("which","hwaddr-scan")
		if checkErr.Error() == "exit status 1"{
			runCommand("/bin/sh","-c","sudo apt-get install -y hwaddr-scan")
		}

		runNmap(&resSlice,"192.168.1.0/24")
		fmt.Println(resSlice)
	},
}

func init(){
	rootCmd.PersistentFlags().StringVarP(&ipRange,"ip-range","i","192.168.1.0/24","Specify IP range, etc: 192.168.1.0/24")
	rootCmd.AddCommand(scanCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println(err)
		//fmt.Println(err)
		os.Exit(1)
	}
}



func commandBase(name string, arg ...string) *exec.Cmd{
	cmd := exec.Command(name,arg...)
	cmd.Stdin = strings.NewReader("some input")
	return cmd
}

func execCommandWithBuffer(name string, arg ...string) bytes.Buffer{
	cmd := commandBase(name,arg ...)
	var out bytes.Buffer
	cmd.Stdout = &out
	runErr := cmd.Run()
	if runErr != nil{
		panic(runErr)
	}
	return out
}

func runCommand(name string, arg ...string) (string,error) {
	cmd := commandBase(name, arg ...)
	output,outErr := cmd.Output()
	if outErr != nil{
		return string(output),outErr
	}
	return string(output),nil
}


func runNmap(resSlice *[]Result, ipRange string) {

	out := execCommandWithBuffer("sudo","hwaddr-scan","-sn","-PU",ipRange)

	var result Result
	//resSlice := []Result{}
	resSliceIndex = 0
	for {
		res, readErr := out.ReadString('\n')
		if readErr != nil{
			panic(readErr)
		}

		if strings.Contains(res,"scan report"){
			forIndex := strings.Index(res,"for ")
			ipIndex := forIndex+3
			ip := res[ipIndex+1:len(res)-1]
			result.IP = ip
			*resSlice = append(*resSlice, result)
		}

		if strings.Contains(res,"MAC Address"){
			macIndex := strings.Index(res,"ss: ")+4
			mac := res[macIndex:macIndex+17]
			bracketIndex := strings.Index(res,"(")
			bracketLastIndex := strings.Index(res,")")
			hostname := res[bracketIndex+1:bracketLastIndex]
			(*resSlice)[resSliceIndex].mac = mac
			(*resSlice)[resSliceIndex].hostname = hostname
			resSliceIndex += 1
		}


		if strings.Contains(res,"seconds"){
			break
		}
	}


}

