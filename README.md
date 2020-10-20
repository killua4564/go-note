## go-note

### Deployment
`docker-compose up -d`

### API Document
#### Account
* create account
	* POST `/api/account`
	* header
		* nil
	* body
		* username
		* password
	* resp
		* success message

* login account
	* POST `/api/account/login`
	* header
		* nil
	* body
		* username
		* password
	* resp
		* token

#### Note
* create note
	* POST `/api/note`
	* header
		* token
	* body
		* topic
		* content (optional)
	* resp
		* sid

* create viewer/owner
	* POST `/api/note/{sid}`
	* header
		* token
	* body
		* is_owner
		* username
	* resp
		* success message

* view note
	* GET `/api/note/{sid}`
	* header
		* token
	* body
		* nil
	* resp
		* topic
		* content

* list note
	* GET `/api/note`
	* header
		* token
	* body
		* nil
	* resp
		* owner
			* sid
			* topic
		* viewer
			* sid
			* topic

* update note
	* PUT `/api/note/{sid}`
	* header
		* token
	* body
		* topic
		* content
	* resp
		* success message

* remove note
	* DELETE `/api/note/{sid}`
	* header
		* token
	* body
		* nil
	* resp
		* success message


### TODO
* database rw split
* refresh note-token
* remove account
* remove account note
* response modulize
* unit testing
