# Apricot

Internal Software for Candidate and Interview Management for Improwised Technologies.

# Development

## Dependencies
* Go 1.6
* [Glide](https://github.com/Masterminds/glide) - Package manager for Go
* PostgreSQL 9.4.4+

## Installation
* Setup go workspace by following this [Guide](https://golang.org/doc/code.html#Organization)
* Set Go vendor flag. Add the following to your `.bashrc` or `.zshrc` along with `$GOPATH` above
```
export GO15VENDOREXPERIMENT=1
```
* Install dependencies
```
$ glide install
```
* Create Database
```
$ create database iims
```
* Database Migration
```
$ go run migration.go
```
* Change Database User
```
$ Change User in Config File
```

## Client Side

* Run Application
```
$ go run index.go
```
* Open
```
$ http://localhost:8000/
```

## Admin Side

* Run Application
```
$ go run admin.go
```
* Open
```
$ http://localhost:8000/
```

## Structure

TODO: Project structure and organization

## Contributing

TODO: Contributing guidelines and code styles

## Testing

TODO: How to run tests
* Admin Side
```
	1. Add Questions.
	2. Add Chellenges.
	3. Add TestCases For Chellenge.
```
* Client Side
```
	1.Enter Email.
	2.Visit link got in Email Id.
```


# Copyright & License

Copyright 2016 Improwised Technologies Pvt Ltd
www.improwised.com

This work is the sole property of Improwised Technologies Pvt Ltd brand
you may not use this work except in compliance with the License.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
