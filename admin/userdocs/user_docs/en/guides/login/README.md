# The login screen

This is where you enter the admin. It appears when you open the admin without
being signed in, and also when your session ends.

## The fields

- **Email** and **Password**: those of your account. The eye icon in the
  password shows and hides what you typed.
- **Keep me signed in on this device**: see below.
- **Sign In**.
- **Forgot your password?**: sends your e-mail a link to create a new
  password.
- **The language**, at the bottom of the screen: the language of the admin.

If the e-mail or the password is wrong, the screen tells you. After several
wrong attempts in a row, the account is locked, and an administrator has to
unlock it. The screen also has a protection against robots: a check the
browser does by itself, or the reCAPTCHA, when the site uses it.

## Keep me signed in on this device

**With the box unchecked**, your session lasts while you use the admin. Each
time you use it, it is renewed. After **an hour unused** it ends, and you have
to sign in again.

**With the box checked**, the session lasts **days** instead of an hour (in
HermonCMS, 30 days, or what the administrator set). It is also renewed while
you use it, so if you use the admin at least once in that time you do not
have to sign in again, not even after you close and open the browser.

It applies **only to this browser, on this device**. On another computer, or
in another browser, you sign in again and choose again.

**Do not check the box** on a shared or public computer. Whoever uses the
computer after you enters the admin with your account.

## When the session ends

- **With a page open:** if the session ends while you are working (a form half
  filled in, say), the login opens in a window over the page. Once you sign
  in, what you were doing goes on, and nothing you typed is lost.
- **From the user menu**, you end it yourself:
  - **Log out** (or **Log out of this account**, with more than one account in
    the browser) ends the session of this browser;
  - **Login Sessions** shows where you are signed in and until when;
  - **Sign out all other sessions** ends all the others, including the ones
    kept signed in on another device.
- **Changing the password** ends the sessions on every device, including the
  ones kept signed in.
