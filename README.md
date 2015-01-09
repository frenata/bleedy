blog
====

go blog project

TODO components:
static page generator
http server
comment handler

static generator
-posts/pages dir in some kind of markdown format
-templates dir 
-function that generates static files 
--routine should check associated comments file based on post id and add them to the page

http server
-simple fileserver that handles requests for static files
-goroutine that checks for new/modified post/pages
--checks for new/changed files and regenerates just that file
-separate handler that checks for POST requests on comment forms and calls comment handler

comment handler
-when requested by server, takes the comment form results and appends them to a file
-send request to regenerator function to regenerate associated page/post 
