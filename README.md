# robin


### Features

* Fiber
* Cron Golang job scheduling for humans. It is inspired by [schedule](<https://github.com/dbader/schedule>).
  


Usage
================

### Quick Start

1.Instsall
~~~
  go get github.com/jiansoft/robin
~~~

2.Use examples
~~~ golang
func main() {
    minute := 11
	second := 50
	//Every friday do once at 11:50:00(HH:mm:ss).
	robin.EveryFriday().At(14, minute, second).Do(runCron, "Friday")

	//Every N day do once at 11:50:00(HH:mm:ss)
	robin.Every(1).Days().At(14, minute, second).Do(runCron, "Days")

	//Every N hours do once at N:50:00(HH:mm:ss).
	robin.Every(1).Hours().At(0, minute, second).Do(runCron, "Every 1 Hours")

	//Every N minutes do once.
	robin.Every(1).Minutes().At(0, 0, second).Do(runCron, "Every 1 Minutes")

	//Every N seconds do once.
	robin.Every(10).Seconds().Do(runCron, "Every 10 Seconds")

}
~~~

## License

Copyright (c) 2017

Released under the MIT license:

- [www.opensource.org/licenses/MIT](http://www.opensource.org/licenses/MIT)