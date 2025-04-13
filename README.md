# Vertias

<p align="center">
  <img src="public/logo.png" height=500/>
</p>

Veritas is a fully customisable JSON-based rule engine implemented using the Go programming language. Existing rule engines exist for Golang applications such as [grule-rule-engine](https://github.com/hyperjumptech/grule-rule-engine) from hyperjumptech. However, the aforementioned engine lacked the required rule customzation and execution flexibility that I desired. 

This project is my attempt at creating an engine that is fully customizable. The engine in this project:
- Utilizes a standardised rule definition structure
- Allows for rules to be evaluated against any user input
- Rules can be designed to specify any kind of response (email response, log response, etc.)

## Technologies

- Go
