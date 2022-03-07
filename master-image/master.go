package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//Copia uma imagem para um nó
func copyFileToNode(filename, dest string, wg *sync.WaitGroup, channel chan int, nodeNumber int) {

	defer wg.Done()

	//sshpass -p '123' scp -o StrictHostKeyChecking=no ./masterInput/image_3.jpg vagrant@172.42.42.103:/home/vagrant/workerInput

	arg0 := "sshpass"
	arg1 := "-p"
	arg2 := "123"
	arg3 := "scp"
	arg4 := "-o"
	arg5 := "StrictHostKeyChecking=no"
	arg6 := filename
	arg7 := "root@" + dest + ":/home/vagrant/workerInput"

	cmd := exec.Command(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)

	//do-while in golang :'(
	for {
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		} else {
			//err = cmd.Wait()
			channel <- nodeNumber
			break
		}
		//err = cmd.Wait()
	}

}

func workerApp(dest string, wgJava *sync.WaitGroup, channelJava chan int, nodeNumber int) {
	defer wgJava.Done()

	arg0 := "sshpass"
	arg1 := "-p"
	arg2 := "123"

	arg3 := "/usr/bin/ssh"
	//arg4 := "-o"
	//arg5 := "StrictHostKeyChecking=no"
	arg6 := "root@" + dest
	arg7 := "/home/vagrant/executeWorkerApp.sh"

	cmd := exec.Command(arg0, arg1, arg2, arg3, arg6, arg7)
	//fmt.Println(cmd.String())

	errCmd := cmd.Run()

	if errCmd != nil {
		fmt.Println(errCmd)
	}
	channelJava <- nodeNumber
	fmt.Printf("WorkerApp.jar do worker%d finalizado!\n\n", nodeNumber)

	//err = cmd.Wait()

}

func workerCopy(dest string, wgCopy *sync.WaitGroup) {
	defer wgCopy.Done()

	arg0 := "sshpass"
	arg1 := "-p"
	arg2 := "123"

	arg3 := "/usr/bin/ssh"
	//arg4 := "-o"
	//arg5 := "StrictHostKeyChecking=no"
	arg6 := "root@" + dest
	arg7 := "/home/vagrant/executeWorkerCopy.sh"

	cmd := exec.Command(arg0, arg1, arg2, arg3, arg6, arg7)
	//fmt.Println(cmd.String())

	errCmd := cmd.Run()

	if errCmd != nil {
		fmt.Println(errCmd)
	}

	//err = cmd.Wait()

}

func clearWorker(dest string, wgCopy *sync.WaitGroup) {
	defer wgCopy.Done()

	arg0 := "sshpass"
	arg1 := "-p"
	arg2 := "123"

	arg3 := "/usr/bin/ssh"
	//arg4 := "-o"
	//arg5 := "StrictHostKeyChecking=no"
	arg6 := "root@" + dest
	arg7 := "/home/vagrant/clearWorker.sh"

	cmd := exec.Command(arg0, arg1, arg2, arg3, arg6, arg7)
	//fmt.Println(cmd.String())

	errCmd := cmd.Run()

	if errCmd != nil {
		fmt.Println(errCmd)
	}

	//err = cmd.Wait()

}

func main() {
	now := time.Now()
	//COLOCAR A MESMA QUANTIDADE DE WORKERS DO VAGRANTFILE
	qtdWorkers := 3

	channel := make(chan int, 6)

	for i := 1; i <= qtdWorkers; i++ {
		channel <- i
	}

	var wg, wgJava, wgCopy, wgClear sync.WaitGroup

	//read images
	files, err := ioutil.ReadDir("./masterInput")
	if err != nil {
		log.Fatal(err)
	}

	var nodeNumber int

	//exercutar as goroutines de acordo com o buffer channel
	for _, f := range files {
		wg.Add(1)

		filename := "./masterInput/" + f.Name()

		nodeNumber = <-channel

		dest := "172.42.42.10" + strconv.Itoa(nodeNumber)

		fmt.Println(filename + " -> " + dest)
		go copyFileToNode(filename, dest, &wg, channel, nodeNumber)

	}

	wg.Wait()
	fmt.Printf("Numero de goroutines: %d\n", runtime.NumGoroutine())

	channelJava := make(chan int, 6)

	for i := 1; i <= qtdWorkers; i++ {
		channelJava <- i
	}

	//executar workerApp.jar em cada Worker
	for {
		wgJava.Add(1)
		i := <-channelJava
		dest := "172.42.42.10" + strconv.Itoa(i)

		fmt.Printf("Executando workerApp.jar no worker%d\n\n", i)
		go workerApp(dest, &wgJava, channelJava, i)
		fmt.Printf("Numero de goroutines: %d\n", runtime.NumGoroutine())

		if i == qtdWorkers {
			break
		}
	}
	fmt.Println("Aguardando término de processamento ...")
	wgJava.Wait()
	fmt.Printf("Numero de goroutines: %d\n", runtime.NumGoroutine())

	//executar workerCopy.jar em cada Worker
	i := 1
	for {
		wgCopy.Add(1)

		dest := "172.42.42.10" + strconv.Itoa(i)

		fmt.Printf("Copiando arquivos de worker%d\n\n", i)
		go workerCopy(dest, &wgCopy)

		wgCopy.Wait()
		i++
		if i == qtdWorkers+1 {
			break
		}
	}

	//executar clearWorker.sh em cada Worker
	i = 1
	for {
		wgClear.Add(1)

		dest := "172.42.42.10" + strconv.Itoa(i)

		fmt.Printf("Limpando arquivos de worker%d\n\n", i)
		go clearWorker(dest, &wgClear)

		wgClear.Wait()
		i++
		if i == qtdWorkers+1 {
			break
		}
	}

	defer func() {
		fmt.Println("RELATÓRIO")
		fmt.Println("Quantidade de workers: ", qtdWorkers)
		fmt.Println("Quantidade de imagens: ", len(files))
		fmt.Println("Tempo de execução: ", time.Since(now))
	}()
}
