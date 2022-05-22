fluroclock-serial-controller
====================

API for controlling the Fluroclock. 

API Endpoints
=============


/panel/numeric
--------------
* Method: POST
* Accepts: application/json
* Body:
```json
{
    "panel_id": "1",
    "value": "0"
}
```
/panel/direct
-------------
* Method: POST
* Accepts: application/json
* Body:
```json
{
    "panel_id": "2",
    "a": "0",
    "b": "0",
    "c": "1",
    "d": "0",
    "e": "1",
    "f": "1",
    "g": "1"
}
```


/panel/alpha
--------------
* Method: POST
* Accepts: application/json
* Body:
```json
{
    "panel_id": "4",
    "alpha": "t"
}
```

/panel/status
--------------
* Method: GET
* Returns:
```json
[
    {
        "Port": {},
        "PortPath": "/dev/ttyUSB50",
        "PanelID": "1",
        "State": {
            "mode": "alpha",
            "bcd_data": null,
            "alpha_data": {
                "panel_id": "1",
                "alpha": "H"
            },
            "direct_data": null
        }
    },
    {
        "Port": {},
        "PortPath": "/dev/ttyUSB51",
        "PanelID": "2",
        "State": {
            "mode": "direct",
            "bcd_data": null,
            "alpha_data": null,
            "direct_data": {
                "panel_id": "2",
                "a": "1",
                "b": "0",
                "c": "1",
                "d": "1",
                "e": "1",
                "f": "1",
                "g": "0"
            }
        }
    },
    {
        "Port": {},
        "PortPath": "/dev/ttyUSB52",
        "PanelID": "3",
        "State": {
            "mode": "alpha",
            "bcd_data": null,
            "alpha_data": {
                "panel_id": "3",
                "alpha": "L"
            },
            "direct_data": null
        }
    },
    {
        "Port": {},
        "PortPath": "/dev/ttyUSB53",
        "PanelID": "4",
        "State": {
            "mode": "alpha",
            "bcd_data": null,
            "alpha_data": {
                "panel_id": "4",
                "alpha": "P"
            },
            "direct_data": null
        }
    }
]
```


/clockupdate
--------------
* Method: POST
Triggers a refresh of the display to show the current time.  Called by fluroclock-clockupdater every minute if enabled.

/clock/disable
--------------
* Method: POST
* Returns: 
`"clock mode disabled"`

Disables the triggering of fluroclock-clockupdater by deleting a file from /etc/ that the timer job requires to run.

/clock/enable
--------------
* Method: POST
* Returns: 
`"clock mode enabled"`

Enables the triggering of fluroclock-clockupdater by writing a file into /etc/ for the timer job to read.  If the file doesn't exist, the clock update mode doesn't occur. 

/clock/isenabled
--------------
* Method: GET
* Returns:
```json
{
    "clock_mode": true
}
```
or 
```json
{
    "clock_mode": false
}
```


/tzdata
--------------
* Method: GET
* Returns:
```json
[
    "Africa/Abidjan",
    "Africa/Accra",
    "Africa/Addis_Ababa",
    "Africa/Algiers",
    "Africa/Asmara",
    "Africa/Asmera",
    "Africa/Bamako",
    "Africa/Bangui",
    "Africa/Banjul",
    ....
]
```

