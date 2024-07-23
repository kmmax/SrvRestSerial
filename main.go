/*
	Пример HTTP-запросов
	Запись в порт массива байт 0..8 (9 байт)
		curl localhost:8080/write?data=012345678
	Чтение из в порта
		curl localhost:8080/read
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tarm/serial"
)

func main() {
	// Настройка последовательного порта
	c := &serial.Config{
		Name:        "COM5", // Замените на Ваш порт
		Baud:        19200,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
		ReadTimeout: time.Second * 5,
	}

	// Открытие последовательного порта
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Ошибка открытия последовательного порта: %v", err)
	}
	fmt.Printf("%s is opened\n", c.Name)
	defer port.Close()

	// Обработчик для чтения данных из порта
	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 128)
		n, err := port.Read(buf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(buf[:n])
		fmt.Fprintf(w, "%s", buf[:n])
	})

	// Обработчик для записи данных в порт
	http.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {
		data := []byte(r.FormValue("data"))
		fmt.Println(data)
		_, err := port.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Данные записаны в последовательный порт")
	})

	log.Println("Запуск HTTP-сервера на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
