# URL shortener
My own iplementation of URL shortener

### Dependencies:
* Go 1.19
* GORM
* SQLite
* Chi

During design Chi router was selected for the project because:
1. I do not need support for custom routing rules or route reversing
2. I do not need support for variables in URL paths
3. I do not need support got host based routing
4. I do not need to support for conflicting routes or regexp route patterns
5. It is not important for me to handle OPTIONS request or "Allow" header on 405 response
6. Will learn a library which is offered like included in micro-service skeleton repository