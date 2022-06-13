package mqttTask

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"lib/models"
)
var client mqtt.Client

func SendAppointment(times int){
	appointmentList := models.GetNowAppointmentList(times)
	for _,appointment := range appointmentList{
		seat := models.GetSeat(appointment.SeatID)
		student,_ := models.GetStudent(appointment.StudentID)
		student.Password = ""
		topic := "device/"+ seat.SeatID
		text,_:= json.Marshal(student)
		token := client.Publish(topic, 0, false, text)
		token.Wait()
	}
}

func SendUndoneAppointment(times int){
	appointmentList := models.GetNowUndoneAppointmentList(times)
	for _,appointment := range appointmentList{
		seat := models.GetSeat(appointment.SeatID)
		topic := "device/"+ seat.SeatID
		text := "undone"
		token := client.Publish(topic, 0, false, text)
		token.Wait()
	}
}

func MqttInitialization(){
	var broker = "127.0.0.1"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	//opts.SetUsername("emqx")
	//opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
