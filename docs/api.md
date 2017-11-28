FORMAT: 1A

# Moppi API
Moppi API is 🌟 Universe for Mesos, Marathon and Chronos in many different KVs.

## Authentication
Currently the Moppi API does not provide authenticated access.

## Error States
The common [HTTP Response Status Codes](https://github.com/for-GET/know-your-http-well/blob/master/status-codes.md) are used.

# Moppi API Root [/]
Moppi API entry point.

# Group Universes

## Universes [/universes]
A single object containing all available universes.

### Create a universe [POST]

Creates a new universe in KV.

+ Request Create Universe (application/json)
    {
        name: "Development",
        version: "1",
        description: "Contains all the packages in developmen"
    }

+ Response 200

### List all universes [GET]

+ Response 200 (application/json)
    [
        {
            "description": "Contains all the packages in development",
            "name": "Development",
            "version": "1",
            "href": "/moppi/universes/dev"
        },
        {
            "description": "Testing",
            "name": "Testing",
            "version": "1",
            "href": "/moppi/universes/testing"
        }
    ]

## Universe [/universes/{universe}]

A Universe object has the following attributes:

+ description
+ name
+ version
+ href

+ Parameters
    + universe: 1 (required, string) - Name of the universe in form of a string

### Delete a universe [DELETE]

+ Response 200

### Get a universe [GET]

+ Response 200 (application/json)

    {
        "description": "Contains all the packages in development",
        "name": "Development",
        "version": "1",
        "href": "/moppi/universes/dev"
    }

### List all Packages [GET /packages]

A list containing all packages available in a universe.

+ Response 200 (application/json)

    [
        "example"
    ]

### Create new Package or new Revision [POST /packages/{package}]

Creates a new Package if there is no package, or creates a new revision if though a package already exists.

+ Parameters
    + package:1 (required, string) - Name of the package in form of a string

+ Request Create a new package or revision (application/json)

    {
        "marathon": [],
        "chronos": [],
        "install": {},
        "uninstall": {}
    }

### List all Revisions [GET /packages/{package}]

A list of all package revisions.

+ Parameters
    + package: 1 (required, string) - Name of the package in form of a string

+ Response 200 (application/json)

    [
        "1"
    ]

### Get a Package [GET /packages/{package}/{revision}]

An package object contains a package.

+ Parameters
    + revision: 1 (required, int) - Revision of the package in form of an integer

+ Response 200 (application/json)

    {
        chronos: [],
        marathon: []
        install: {},
        uninstall: {}
    }

### Delete a Package [DELETE /packages/{package}]

Returns nothing

+ Parameters
    + revision: 1 (required, int) - Revision of the package in form of an integer

+ Response 201 (application/json)

    {
        chronos: [],
        marathon: []
        install: {},
        uninstall: {}
    }

# Group Install/Uninstall 

### Install a package [POST /install]

Triggers the installation of a new package

+ Response 201

+ Request Install a package (application/json)

    {
        "universe": "dev",
        "revision": 1,
        "name": "example",
        "config": {}
    }

### Uninstall a package [POST /uninstall]

Triggers the installation of a new package

+ Response 201

+ Request Uninstall a package (application/json)

    {
        "universe": "dev",
        "revision": 1,
        "name": "example",
        "config": {}
    }
