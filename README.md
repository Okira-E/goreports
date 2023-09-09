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

## OpenAPI Documentation

GoReports has an OpenAPI documentation that can be found at `/swagger` endpoint when the server is running.

## Usage and Documentation

GoReports is a command line tool. You can run `goreports --help` to see the available commands and flags.

### Start GoReports using Docker

1. Build the Docker image using the following command:

```shell
docker build -t goreports \
  --build-arg db_dialect=your_db_dialect \
  --build-arg db_user=your_db_username \
  --build-arg db_password=your_db_password \
  --build-arg db_host=your_db_host \
  --build-arg db_port=your_db_port \
  --build-arg db_name=your_db_name .
```

Note: Database information is required to execute the queries in the templates and fetch the data for the reports.

2. Run the Docker image using the following command:

```shell
docker run -p 3200:3200 goreports
```

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
This is a snippet of a template that uses all the syntaxes mentioned above:
```html
<div style="max-width: 800px; margin: 0 auto;">
    <h1>Payment History Report</h1>

    <p>[P[extra_param]]</p> <!-- This is a parameter passed to the template at request/render time -->

    <table>
        <tr>
            <th>Date</th>
            <th>Invoice Number</th>
            <th>Amount Paid</th>
        </tr>
        {{#each [Q[SELECT creation_date, invoice_number, amount_paid FROM payments WHERE customer_id =
        [P[customer_id]] ]]}} <!-- This is a query passed to the template with a parameter in it -->
        <tr>
            <td>{{creation_date}}</td>
            <td>{{invoice_number}}</td>
            <td>{{amount_paid}}</td>
            {{/each}}
        <tr>
    </table>
</div>
```
It generates (along with hidden CSS code) this [report](./examples/payment_history.pdf) when rendered. As you can see, The SQL query is executed with the `customer_id` param (passed at each render request in the API) and its multiple results are passed to the handlebars loop.


### Save a report

To save a report, send a POST request to GoReports' server at `/report/save` endpoint with the following JSON body:

```json
{
  "name": "required",
  "title": "required",
  "description": "optional",
  "header": "<html>optional</html>",
  "body": "<html>required</html>",
  "footer": "<html>optional</html>"
}
```

The `name` field is the name of the report. Rendering the report will require this name.

The `body` field is the template to be parsed and rendered into PDF.

`header` and `footer` fields are optional and will be prepended and appended to the `body` respectively on each page.

***Note: page numbers are generated at render time and replace the footer or the header if aer positioned at the bottom or the top of the page respectively.***

### Render a report

After saving a report you can render it by sending a POST request to `/report/render` endpoint with the following JSON body and options:

```json
{
  "reportName": "payment_history",
  "printingOptions": {
    "paperSize": "A4",
    "landscape": false,
    "marginTop": 20,
    "marginRight": 0,
    "marginBottom": 10,
    "marginLeft": 0,
    "pageNumbers": {
      "enabled": false,
      "position": "bottom-center" // top-left, top-center, top-right, bottom-left, bottom-center, bottom-right
    }
  },
  "params": {
    "customer_id": 2,
    "extra_param": "This is a parameter passed from the client."
  }
}
```
***Note: It is on you to provide the necessary margins to show headers, footers and page numbers. Page numbers replace footers/headers if positioned on the bottom/top.***

The report will be rendered into PDF and sent as a buffer in the response.

### Delete a report

To delete a report, send a DELETE request to GoReports' server at `/report/delete` endpoint with the following JSON body:

```json
{
  "reportName": "required"
}
```

### List all reports

To list all reports, send a GET request to GoReports' server at `/report/list` endpoint.

## Contributing

Pull requests are always welcomed and encouraged. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/). See the [LICENSE](LICENSE) file for details.