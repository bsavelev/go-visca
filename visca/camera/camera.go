package camera

import (
    "github.com/pkg/errors"
    "log"
)

// #include "libvisca.h"
import "C"

var ErrRequestCamera = errors.New("cam error request")

type CameraConfig struct {
    Port			string
    Rate			int
}

type Camera struct {
    iface      C.VISCAInterface_t
    camera     C.VISCACamera_t
    cameraNum	C.int
    currentZoom float64
    portName	string
    maxZoom 	float64
}

func (cam *Camera) Reconnect(cfg *CameraConfig) error {
    cam.portName = cfg.Port

    if err := cam.close(); err != nil {
	log.Println("Camera.Reconnect: close failed")
    }

    if err := cam.openSerial(cfg.Port, cfg.Rate); err != nil {
        log.Println("Camera.Reconnect: openSerial failed")
	return err
    }
    log.Print("openSerial success")

    cam.iface.broadcast = 0

    if err := cam.setAddress(); err != nil {
	log.Println("Camera.Reconnect: setAddress failed")
    }

    cam.camera.address = 1

    if err := cam.clear(); err != nil {
	log.Println("Camera.Reconnect: clear failed")
    }

    if err := cam.getCameraInfo(); err != nil {
	log.Println("Camera.Reconnect: getCameraInfo failed")
	return err
    }
    log.Print("Vendor: ", cam.camera.vendor);
    log.Print("Model: ", cam.camera.model);
    log.Print("ROM version: ", cam.camera.rom_version);
    log.Print("Socket num: ", cam.camera.socket_num);

    if err := cam.disableZoom(); err != nil {
	log.Println("Camera.Reconnect: disableZoom failed")
	return err
    }
    cam.maxZoom = 32768
    return nil
}

func (cam *Camera) StopCameraZoom() error {

    return  cam.setZoomStop()
}

func (cam *Camera) SetCameraZoom(val uint) error {
    return cam.setZoom(val)
}

func (cam *Camera) SetMaxZoom(val float64) {
    cam.maxZoom = val
}

func (cam *Camera) ZoomIn () error {
    if err := cam.setZoomTele(); err != nil {
	return err
    }
    if err := cam.setZoomTeleWithSpeed(); err != nil{
	return err
    }
    return nil
}

func (cam *Camera) ZoomOut() error {
    log.Println("[visca] zoom out")

    if err := cam.setZoomWide(); err != nil {
	return err
    }
    if err := cam.setZoomWideWithSpeed(); err != nil{
	return err
    }
    return nil
}

func (cam *Camera) request (res C.uint) (err error) {
    if int(res) == C.VISCA_SUCCESS {
	return nil
    }
    return ErrRequestCamera
}

func (cam *Camera) setZoomTele() error{
    return cam.request(C.VISCA_set_zoom_tele(&cam.iface, &cam.camera))
}

func (cam *Camera) setZoomTeleWithSpeed() error{
    return cam.request(C.VISCA_set_zoom_tele_speed(&cam.iface, &cam.camera, 7))
}

func (cam *Camera) setZoomWide() error{
    return cam.request(C.VISCA_set_zoom_wide(&cam.iface, &cam.camera))
}

func (cam *Camera) setZoomWideWithSpeed() error{
    return cam.request(C.VISCA_set_zoom_wide_speed(&cam.iface, &cam.camera, C.uint32_t(7)))
}

func (cam *Camera) setZoomStop() error {
    return cam.request(C.VISCA_set_zoom_stop(&cam.iface, &cam.camera))
}

func (cam *Camera) restoreZoom() (error) {
    return cam.request(C.VISCA_set_zoom_value(&cam.iface, &cam.camera, C.uint32_t(0)))
}

func (cam *Camera) openSerial(serialPort string, rate int) (err error) {
    return cam.request(C.VISCA_open_serial(&cam.iface, C.CString(serialPort)))
}

func (cam *Camera) setAddress() error {
    return cam.request(C.VISCA_set_address(&cam.iface, &cam.cameraNum))
}

func (cam *Camera) close ()(error) {
    return cam.request(C.VISCA_close_serial(&cam.iface))
}

func (cam *Camera) getCameraInfo () (err error) {
    err = cam.request(C.VISCA_get_camera_info(&cam.iface, &cam.camera))
    if err != nil {
	return err
    }
    return nil
}

func (cam *Camera) clear() (err error) {
    return cam.request(C.VISCA_clear(&cam.iface, &cam.camera))
}

func (cam *Camera) disableZoom () (err error) {
    cam.request(C.VISCA_set_dzoom_mode(&cam.iface, &cam.camera,  C.VISCA_DZOOM_OFF))
    cam.request(C.VISCA_set_dzoom_mode(&cam.iface, &cam.camera,  C.VISCA_DZOOM_OFF))
    return nil
}

func (cam *Camera) getZoom () (zoom float64, err error) {
    var val C.uint16_t
    if err = cam.request(C.VISCA_get_zoom_value(&cam.iface, &cam.camera, &val)); err != nil {
	return
    }
    log.Println(`from visca`, val)
    zoom = float64(val)

    return zoom, nil
}

func (cam *Camera) setZoom(val uint) (err error) {
    return cam.request(C.VISCA_set_zoom_value(&cam.iface, &cam.camera, C.uint32_t(val)))
}


func (cam *Camera) Check() (err error) {
    if err = cam.setZoom(0); err != nil {
	log.Print("error setZoom: ", 0)
	return err
    }
    if err = cam.setZoom(uint(cam.maxZoom)); err != nil {
	log.Print("error setZoom: ", cam.maxZoom)
	return err
    }
    val, err := cam.getZoom()
    if err != nil {
	log.Print("error getZoom")
	return err
    } else {
        log.Print("Zoom: ", val)
    }
    return nil
}
