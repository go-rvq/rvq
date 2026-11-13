# Permission sets

A set is a name for permissions given together. A line of a role uses it
with `%`: `admin:posts:%readers` lists the posts and sees each one. The roles
are in {%= admin.model("roles").link %}.

Here the sets of the admin are made, beside the system's (`readers`,
`authors`, `editors`…):

- **Name**: how the lines write it after `%` — lower case letters, digits
  and `_`, a letter first. It may not be one of a set of the system.
- **Definition**: what it gives by kind of resource, for example
  `{model: {@list, <*>:@get}, singleton: @get}`.
- **Titles** and **Descriptions**: one in each language of the admin. They
  are what the field of the permissions shows.

The line keeps the name. Changing a set changes what every line of it
gives, at each user's next request. A deleted set leaves its lines kept and
ignored, shown red.

All of it, with the sets of the system, in
{%= admin.doc("guides/policies/01-permissions/02-sets").link %}.
