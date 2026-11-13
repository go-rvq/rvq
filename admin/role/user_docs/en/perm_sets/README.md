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
- **Title** and **Description**: in the language of the record. They are
  what the field of the permissions shows.

A set is **localized by the languages of the admin**, as pages and posts by
the site's: a record in each language, each with its title and description.
The listing shows those of the admin's language; **Localize** copies a set
into another language and **Localizations** shows those it has. The name and
the definition are the set's: saved in one language, they hold in all. A
language of the admin made from another gets the sets already translated, by
the automatic translator, when it is on.

The sets **of the system** (`readers`, `authors`…) are registered here by
themselves, marked *Of the system*: their name and definition are the
code's, read only, and they are not deleted. Their title and description may
be translated.

The detail shows who made it and who changed it last, and when. The tab of
the history keeps each change.

The line keeps the name. Changing a set changes what every line of it
gives, at each user's next request. A deleted set leaves its lines kept and
ignored, shown red.

All of it, with the sets of the system, in
{%= admin.doc("guides/policies/01-permissions/02-sets").link %}.
