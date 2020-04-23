
### API reference

```
Endpoint:               Variables:              Comment:
/restart*                                       Graceful server shutdown.
                                                Requires admin rights
                        wipe        string      Wipes database and images if "yes"


/register                                       Register a user account
                        user        string      Username
                        pass        string      Password
                        email       string      Email address


/login                                          Login to get higher access level
                        user        string      Username
                        pass        string      Password


/getusers               <null>                  Retrieves list of all users
                                                in database


/gettags                <null>                  Retrieves list of all tags
                                                in database


/getthemes              <null>                  Retrieves list of all theme
                                                files in static/css/themes/


/getdeck                                        Request for slide deck
                        id          int         Retrieve saved deck with id#
                        amount      int         # of slides requested
                        lang        string      language code 'en', 'sv' etc
                        tags        string      tags on which to base deck


/addtext*                                       Adds new text to database
                        tags        string      Which tags to associate text with
                        ttext       string      Title text
                        btext       string      Body text


/addimg*                                        Adds new images to the database
                        file        file        The image file itself
                        tags        string      Which tags to associate the image with


/chuser*                                        Change user settings
                                                Some ops requires admin rights
                        tuser       string      User to edit
                        pass        string      New password (if applicable)
                        op          int         Operation:
                                                    0: Make admin
                                                    1: Remove admin rights
                                                    2: Change password
                                                    3: Remove user account
                                                    4: Email new password


/report*                                        Reports inappropriate content to admin email
                        id          int         Deck number
                        slide       int         Slide to report
                        msg         string      Motivation for report


/remove*                                        Removes object from db
                                                Requires admin rights
                        id          int         Object id
                        type        string      Object type:
                                                    img: image
                                                    ttext: title
                                                    btext: body text


/feedback*                                      Give feedback on user experience
                        msg         string      The feedback info itself

```
Endpoints marked with `*` requires the user to be logged in, authenticated by `user` & `skey` being included with the request.

