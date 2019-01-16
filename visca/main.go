package main

import (
    "hi-it.spb.ru/visca/camera"
    "log"
)

func main() {
    log.Print("Start visca test")
    var cam = &camera.Camera{}
    cam.Reconnect(&camera.CameraConfig{
        Port: "/dev/tty_twiga",
        Rate: 9600,
    })
    cam.Check()
    log.Print("End visca test")
}
