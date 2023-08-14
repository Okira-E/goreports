# GoReports

GoReports is a report generation tool that allows you to build and generate dynamic reports using a custom
handlebars syntax that allows you to add SQL queries directly into the template and pass parameters to the queries.

See the [Example](#example) section for a sample template.

## Where to use GoReports?

Wherever there is a `Print as PDF` button, GoReports can be used to generate the PDF. Some examples are:
- Invoices
- Payment history
- Sales reports

## Prerequisites

You need [wkhtmltopdf](https://wkhtmltopdf.org/downloads.html) installed and in your PATH for GoReports to work.

## Installation

You can install the GoReports binary from the [releases page](https://github.com/Okira-E/goreports/releases).

## Building from source

### Requirements

- Go 1.18 or higher
- Enable `CGO_ENABLED` by `go env -w CGO_ENABLED=1`. This is required for the SQLite driver to work.

### Building locally

```shell
go env -w CGO_ENABLED=1
go build
```

## Usage and Documentation

GoReports is a command line tool. You can run `goreports --help` to see the available commands and flags.

### Set up GoReports on your machine

```shell
goreports init
```

This will ask you for your external database credentials to execute the queries in the templates from and create
a `config.json` storing the credentials, as well as, a `data/` directory in goreports config directory based on your OS.
An internal [SQLite](https://www.sqlite.org/index.html) database will be created in the `data/` directory to store the
reports.

### Start the server

```shell
goreports start
```

This will start goreports server on port 3200.

### Template syntax

GoReports uses an extended handlebars syntax to parse and render templates. The syntax is as follows:

- `[P[PARAMETER_NAME]]` - This will be replaced with the value of the parameter passed to the template.
- `[Q[SQL_QUERY]]` - This will be replaced with the results of the SQL query. If the query returns multiple rows,
  the result will be an array of objects. If the query returns a single row, the result will be a be inserted directly into the HTML.
- `{{#each [Q[SQL_QUERY]]}}` - This will execute the (multiple results) SQL query and pass the result array to the
  handlebars block. You can then access the properties of each object in the array like you would in a normal handlebars.

#### Example
```html
<div style="max-width: 800px; margin: 0 auto;">
    <h1>Payment History Report</h1>

    <p>[P[extra_param]]</p> <!-- This is a parameter passed to the template at render time -->

    <table>
        <tr>
            <th>Date</th>
            <th>Invoice Number</th>
            <th>Amount Paid</th>
        </tr>
        {{#each [Q[SELECT creation_date, invoice_number, amount_paid FROM payments WHERE customer_id =
        [P[customer_id]] ]]}} <!-- This is a query passed to the template with a parameter in it -->
        <tr>
            <td>{{datify creation_date}}</td>
            <td>{{invoice_number}}</td>
            <td>{{amount_paid}}</td>
            {{/each}}
        <tr>
    </table>
</div>
```

### Save a report

To save a report, send a POST request to GoReports' server at `/report/save` endpoint with the following JSON body:

```json
{
  "name": "required",
  "title": "required",
  "description": "optional",
  "template": "<html>required</html>"
}
```

The `name` field is the name of the report. Rendering the report will require this name.

The `template` field is the template to be parsed and rendered into PDF.

### Render a report

After saving a report you can render it by sending a POST request to `/report/render` endpoint with the following JSON body and options:

```json
{
  "reportName": "payment_history",
  "printingOptions": {
    "paperSize": "A4",
    "landscape": false
  },
  "params": {
    "customer_id": 2,
    "extra_param": "This is a parameter passed from the client.",
  }
}
```

The report will be rendered into PDF and sent as a buffer in the response (future versions may allow you to save the PDF to S3 buckets).


## Contributing

Pull requests are always welcomed and encouraged. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/). See the [LICENSE](LICENSE) file for details.