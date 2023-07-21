# GoInvoicer
> Tried to keep the project as simple as possible didn't use any fancy pattern or any framework. You can Implement the frontend as you like with the help of the docs below. Please create a issue if you find any kind of problem or if you want to contribute in this project. I will try to implement a frontend soon with vue.js.
> NOTE: It uses wkhtmltopdf for pdf generation. The app includes both binaries (Windows/Linux) for wkhtmltopdf in ./bin folder. Keep this in mind for your deployment strategy.

> ### [API Documentation](https://documenter.getpostman.com/view/28287667/2s946bBus4) (Postman)

## Local Installation
- `` git clone https://github.com/FahimAnzamDip/GoInvoicer.git ``
- `` cd to/project/path ``
- `` copy .env.example to .env ``
- `` go mod tidy ``
- `` go run main.go ``

## Achiveable with this API

- **Informative dashboard**
    - Displays the Summary of all information in one singular view.

- **Client Management**
    - Adding Clients and creating accounts for individual clients/organizations to keep track of the invoice and payments history.

- **Expense Management**
    - Add all types of expenses to the app and track your expense details.

- **Invoicing**
    - Creating invoices/Generating Bills and send to client via email. You can use the app for making payment against the generated invoices.

- **Recurring Invoicing**
    - Generate and Autosend the recurred invoice to the client (Monthly/Quarterly/Semi-Annually/Annually).

- **Estimates**
    - Create estimates and send it to your client via email. Clone estimates to invoice.

- **Product Management**
    - You are able to add Products/Services according to their category.

- **Reports**
    - See different reports to learn about client and their payment status at a glance.
