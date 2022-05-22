package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	promMW "github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	apiCalls = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_calls",
		Help: "The total number of api calls",
	})
)

func PanelBCDModeHandler(pm *[]PortMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiCalls.Inc()
		bdata := new(BCDData)
		if err := c.Bind(bdata); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		if bdata.PanelID == "" {
			return c.JSON(http.StatusBadRequest, "PanelID is required")
		}
		if bdata.Value == "" {
			return c.JSON(http.StatusBadRequest, "Value is required")
		}
		pm, pme := GetPortMapping(bdata.PanelID, *pm)
		if pme != nil {
			return c.JSON(http.StatusNotFound, "Port Mapping not found")
		}
		if pm.Port == nil {
			return c.JSON(http.StatusNotFound, "Port is not found")
		}
		realValue, cerr := strconv.Atoi(bdata.Value)
		if cerr != nil {
			return c.JSON(http.StatusBadRequest, "Value must be an integer between 0 and 9")
		}
		if realValue < 0 || realValue > 9 {
			return c.JSON(http.StatusBadRequest, "Value must be between 0 and 9")
		}

		code := fmt.Sprintf("+P%sMIV%d-", bdata.PanelID, realValue)
		sse := SendSerialString(code, pm.Port)
		if sse != nil {
			return c.JSON(http.StatusInternalServerError, sse)
		}

		pm.State.Mode = "bcd"
		pm.State.AlphaData = nil
		pm.State.BCDData = bdata
		pm.State.DirectData = nil
		return c.JSON(http.StatusOK, "OK")

	}
}

func PanelDirectModeRequestHandler(pm *[]PortMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiCalls.Inc()
		ddata := new(DirectData)
		if err := c.Bind(ddata); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		if ddata.PanelID == "" {
			return c.JSON(http.StatusBadRequest, "PanelID is required")
		}
		if ddata.A == "" || ddata.B == "" || ddata.C == "" || ddata.D == "" || ddata.E == "" || ddata.F == "" || ddata.G == "" {
			return c.JSON(http.StatusBadRequest, "All values are required")
		}
		if (ddata.A != "0" && ddata.A != "1") ||
			(ddata.B != "0" && ddata.B != "1") ||
			(ddata.C != "0" && ddata.C != "1") ||
			(ddata.D != "0" && ddata.D != "1") ||
			(ddata.E != "0" && ddata.E != "1") ||
			(ddata.F != "0" && ddata.F != "1") ||
			(ddata.G != "0" && ddata.G != "1") {
			return c.JSON(http.StatusBadRequest, "All values must be 0 or 1")
		}

		pm, pme := GetPortMapping(ddata.PanelID, *pm)
		if pme != nil {
			return c.JSON(http.StatusNotFound, "Port Mapping not found")
		}
		if pm.Port == nil {
			return c.JSON(http.StatusNotFound, "Port is not found")
		}
		code := fmt.Sprintf("+P%sMXA%sB%sC%sD%sE%sF%sG%s-", ddata.PanelID, ddata.A, ddata.B, ddata.C, ddata.D, ddata.E, ddata.F, ddata.G)
		sse := SendSerialString(code, pm.Port)
		if sse != nil {
			return c.JSON(http.StatusInternalServerError, sse)
		}
		pm.State.Mode = "direct"
		pm.State.AlphaData = nil
		pm.State.BCDData = nil
		pm.State.DirectData = ddata
		return c.JSON(http.StatusOK, "OK")

	}
}

func PanelAlphaModeRequestHandler(pm *[]PortMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiCalls.Inc()
		adata := new(AlphaData)
		if err := c.Bind(adata); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		if adata.PanelID == "" {
			return c.JSON(http.StatusBadRequest, "PanelID is required")
		}
		if adata.Alpha == "" {
			return c.JSON(http.StatusBadRequest, "Alpha is required")
		}
		pm, pme := GetPortMapping(adata.PanelID, *pm)
		if pme != nil {
			return c.JSON(http.StatusNotFound, "Port Mapping not found")
		}
		if pm.Port == nil {
			return c.JSON(http.StatusNotFound, "Port is not found")
		}
		letter_hx, err := GetLetter(adata.Alpha)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Alpha must be in the list of letters: "+strings.Join(GetValidLetters(), ","))
		}

		ddata := ConvertBitstringToDataStruct(ConvertIntToBitstring(letter_hx))
		code := fmt.Sprintf("+P%sMXA%sB%sC%sD%sE%sF%sG%s-", adata.PanelID, ddata.A, ddata.B, ddata.C, ddata.D, ddata.E, ddata.F, ddata.G)

		sse := SendSerialString(code, pm.Port)
		if sse != nil {
			return c.JSON(http.StatusInternalServerError, sse)
		}
		pm.State.Mode = "alpha"
		pm.State.AlphaData = adata
		pm.State.BCDData = nil
		pm.State.DirectData = nil

		return c.JSON(http.StatusOK, "OK")
	}
}

func PanelStatusHandler(pm *[]PortMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiCalls.Inc()
		ClockModeEnabled = InClockModeFromRestart()
		return c.JSON(http.StatusOK, pm)
	}
}

func ClockModeHandler(pm *[]PortMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiCalls.Inc()
		ClockModeUpdateTime(pm)

		return c.JSON(http.StatusOK, "OK")
	}
}

func TZDataHandler(tzdata *[]string) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiCalls.Inc()
		return c.JSON(http.StatusOK, tzdata)
	}
}

func ClockDisableHandler(pm *[]PortMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiCalls.Inc()
		ClockModeEnabled = true
		fileError := os.Remove("/etc/fclock-clock-enable")
		if fileError != nil {
			log.Fatal(fileError)
		}
		ResetDisplay()

		return c.JSON(http.StatusOK, "clock mode disabled")
	}
}
func ClockEnableHandler(pm *[]PortMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiCalls.Inc()
		ClockModeEnabled = true
		Touch("/etc/fclock-clock-enable") //  this file will be detected by clockupdater and will cause it to run
		ClockModeUpdateTime(pm)
		return c.JSON(http.StatusOK, "clock mode enabled")
	}

}

func QueryClockModeEnabledHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		apiCalls.Inc()
		sstate := SystemState{
			ClockMode: ClockModeEnabled,
		}
		return c.JSON(http.StatusOK, sstate)
	}
}

func StartServer(pm *[]PortMap) {
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		log.Fatal("$PORT must be set")
	}
	tzdata, tzr := LoadTZData("tz.json")
	if tzr != nil {
		log.Fatal(tzr)
	}
	ClockModeEnabled = InClockModeFromRestart()
	log.Println("Clock Mode Enabled:", ClockModeEnabled)
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	p := promMW.NewPrometheus("echo", nil)
	p.Use(e)

	e.POST("/panel/numeric", PanelBCDModeHandler(pm))
	e.POST("/panel/direct", PanelDirectModeRequestHandler(pm))
	e.POST("/panel/alpha", PanelAlphaModeRequestHandler(pm))
	e.GET("/panel/status", PanelStatusHandler(pm))

	e.POST("/clockupdate", ClockModeHandler(pm))
	e.POST("/clock/disable", ClockDisableHandler(pm))
	e.POST("/clock/enable", ClockEnableHandler(pm))
	e.GET("/clock/isenabled", QueryClockModeEnabledHandler())
	e.GET("/tzdata", TZDataHandler(tzdata))

	e.Logger.Fatal(e.Start("0.0.0.0:" + httpPort))

}
