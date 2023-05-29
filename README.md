# Mwkhtmltopdf

A wrapper that separates Odoo from Wkhtmltopdf for high availability purposes.

## Intro

Odoo uses Wkhtmltopdf for generating PDF files. However, Wkhtmltopdf is subject to Odoo worker limits as it's a subprocess of Odoo workers.

This project has 2 components, the client and the server, the client is a binary that extracts wkhtmltopdf CLI arguments created by Odoo, and sends them to the server, the server parses the request from client and generates a PDF, then returns it to client. From Odoo's perspective, the client is just a normal Wkhtmltopdf binary, so there's no need to make changes in the Odoo code.

The advantages are:

1. Wkhtmltopdf is no longer subject to Odoo worker limits and thus can scale infinitely regardless of Odoo's resources and limits
2. Wkhtmltopdf can be replaced in a completely different server, or most cases, in a different Docker container. This helps reducing Odoo Docker base image size, and even allows Odoo to be installed in environment where Wkhtmltopdf is not installable, e. g. Alpine Linux.

## Configuration

### Client

1. `MWKHTMLTOPDF_URL`: Base URL of the Mwkhtmltopdf server, if you're using [odooctl](https://github.com/mplus-oss/odoo) deployment stack this value is set to "http://mwkhtmltopdf-server:2777"


## Installation

Include this line in your Dockerfile:

```Dockerfile
COPY --from=ghcr.io/mplus-oss/mwkhtmltopdf-client /usr/local/bin/wkhtmltopdf /usr/local/bin/wkhtmltopdf
```

## Security

Do not ever expose the server! There's no authentication system which makes the server should be forever run in a private network, otherwise it can lead to complete takeover of the OS/container it is running!
