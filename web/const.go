package web

const GOPHILE_HTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Go WebServer</title>

    <style type="text/css">
        * {
            font-family: 'Helvetica', 'Arial', sans-serif;
        }

        body {
            height: 100%;
            width: 80%;
            margin: 0 auto;
            background-color: #333;
        }

        header {
            position: relative;
            width: 100%;
            height: 50px;
            background-color: #AAA;
            border-radius: 0px 0px 15px 15px;
        }

        #logo {
            margin-left: 10px;
            height: 50px;
        }

        #changePass {
            position: absolute;
            right: 115px;
            top: 14px;
        }

        #accountField {
            margin: 5px;
            text-align: center;
            right: 0;
            top: 0;
            position: absolute;
        }

        #uploadField {
            margin-top: 5px;
            border-radius: 15px;
            text-align: center;
            background-color: #CCC;
            padding: 5px;
        }

        nav {
            margin-top: 5px;
            border-radius: 15px;
            background-color: #AAA;
            text-align: center;
            padding: 2px;
        }

        nav form {
            display: inline-block;
        }
        #ErrorLine {
            text-align: center;
            color: #fff;
            margin-top: 10px;
            margin-bottom: 10px;
        }

        #BackButton {
            background: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAQAAABKfvVzAAAAWklEQVR4Ae3LUQqAMAwE0TmP59F6G8m9W6gQkIQiLv1U3N+dx4dn2Fze6dhUnonOm+dNkciLg10SzysFHKCI52xwgSAqD/BAjMoKGSRycLMFBhCP2gD4wet3AqyEQlKENMZpAAAAAElFTkSuQmCC') no-repeat left center;
        }

        #RefreshButton {
            background: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAQAAABKfvVzAAAAsklEQVR4Ad3QMUoDQRSA4a9yWBCPYJFWFrOR3C+FYJUmiZgiZ7HwEEEWCxE8QzIIedbLMsvU+aYb/scwz1Xq7PWyP703jyY1DmJwvtxO5R/CydpCcqPz7F6Zg/DjQaVOOBXylZWRvbAu5CHGI72wKObjEVlIpXw4Mr4caoRz5bNgKRwrPwY2wrZydZjLLlqV5n6FnTIzL54kjaWNLLxLiu58D3Z2sZNMar36lJ0dbbWu0D+h01pq33dROwAAAABJRU5ErkJggg==') no-repeat left center;
        }

        #folderTitle {
            vertical-align: top;
            display: inline-block;
            margin-top: 7px;
        }

        #folderPath {
            border: 1px solid black;
            padding: 2px;
        }

        #createFolder {
            vertical-align: top;
            margin-top: 6px;
            margin-left: 12px;
        }

        #createFolder input:first-child {
            height: 20px;
        }

        table {
            border-spacing: 0 5px;
            width: 100%;
        }

        tr td {
            background: white;
            border-radius: 0;
            padding: 3px 7px;
        }

        td:first-child {
            border-radius: 18px 0 0 18px;
            min-width: 15px;
        }

        td:last-child {
            border-radius: 0 18px 18px 0;
            min-width: 15px;
        }

        input.clickable {
            width: 28px;
            height: 30px;
            border: none;
        }

        input.clickable:hover {
            cursor: pointer;
            background-color: #AAAAAA;
        }

        .DeleteButton {
            background: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAQAAABKfvVzAAAANklEQVR4AWMYCkCB4QHDfyh8AOThAf9xQzI1UA5INX9Uw6gGnOAVVuUvcGvwZ3iJRbnvIM/yAMG7x37sRHSvAAAAAElFTkSuQmCC') no-repeat left center;
        }

        .DownloadButton {
            background: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAQAAABKfvVzAAAAOUlEQVR4AWMYGuA/Ao5qwAANUEWYsIE0LQ2k2IJQTlgLQjlBLRjKCWtBKCdeS8OgStaYkHQNwwoAAEWCcJenlv+/AAAAAElFTkSuQmCC') no-repeat left center;
        }

        .LinkButton {
            background: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAQAAABKfvVzAAAArUlEQVR4Ae3LMUoDQQCF4e8a4hbRNOJFBE+VIFpZRMtopzmAkhMErcRbLBZZBGVmGzGJNhmw2RnYKuBXvOr9+js214oeHZXdPy1NXGl8lCRzS/ug0niQ1ZrYuhZkbYxsjWyKgrS/QdqdCP6DlXNcWJUGzxqX3i3ywdoYQ0++LBxgbK3Dmzt/zdQ63Ar2SCrRVIeh1ktKKq+CgU6nouDemZkoOJF16EbtW21qoK8f5VpvkEyGyc0AAAAASUVORK5CYII=') no-repeat left center;
        }
    </style>
</head>

<body>
<!-- Header with the logo, possibility to change the password and the logout button -->
<header>
    <img id="logo"
         src='data:image/png;base64,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'
         alt="gophile">
    <div id="changePass">
        <form action="/changePassword" method="POST">
            <span>Change password: </span>
            <input type="password" placeholder="Current password" name="currentPass" required>
            <input type="password" placeholder="New password" name="newPass" required>
            <input type="submit" value="Change">
        </form>
    </div>
    <div id="accountField">Hello {{.User}}!
        <form action="/logout" method="POST">
            <input type="submit" value="Logout">
        </form>
    </div>
</header>

<!-- Uploading new files -->
<div id="uploadField">
    <form enctype="multipart/form-data" method="POST" action="/gophile">
        <input name="UploadFile" type="file" required>
        <input type="submit" value="Upload">
    </form>
</div>

<!-- Navigation bar with 'up' navigation, refresh, display of current path, button for creating new folders -->
<nav>
    <form method="post" action="/gophile">
        <input type="hidden" name="Back" value="true">
        <input id="BackButton" class="clickable" type="submit" value="">
    </form>
    <form method="post" action="/gophile">
        <input type="hidden" name="Refresh" value="true">
        <input id="RefreshButton" class="clickable" type="submit" value="">
    </form>
    <span id="folderTitle">Current Path: <span id="folderPath">/{{.CurrentPath}}</span></span>
    <form method="POST" id="createFolder" action="/gophile">
        <label for="FolderName">New Folder</label>
        <input id="FolderName" name="FolderName" type="text" placeholder="Folder name" required/>
        <input type="submit" value="Create"/>
    </form>
</nav>
<!-- for showing hints and error messages -->
<div id="ErrorLine">{{.ActionText}}</div>
<!-- Table for files and folders. Every line shows a single file/folder with its size, date and buttons for download, link and delete -->
<table>
{{range .FileTable}}
    <tr>
    {{if .IsDir}}
        <td>
            <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAQAAABKfvVzAAAAPUlEQVR4AWMYNiCQ4SXDfyTYQEA9RDkpWv4TAV8x+KNqIAxfkqrh/5DXMKrhFVHKXyA0+DO8JEK573DJzQAtzyV302eonwAAAABJRU5ErkJggg=="/>
        </td>
        <td>
            <form class="buttonForm" method="post" action="/gophile">
                <input type="hidden" name="BrowseFolder" value="{{.FileName}}"/>
                <input type="submit" value="{{.FileName}}"/>
            </form>
        </td>
    {{else}}
        <td></td>
        <td>{{.FileName}}</td>
    {{end}}
        <td>{{.FileSize}}</td>
        <td>{{.FileDate}}</td>
    {{if eq .IsDir false}}
        <td>
            <form method="post" action="/gophile">
                <input type="hidden" name="DownloadFile" value="{{.FileName}}"/>
                <input class="DownloadButton clickable" value="" type="submit"/>
            </form>
        </td>
        <td class="clickable">
            <form method="post" action="/gophile">
                <input type="hidden" name="LinkFile" value="{{.FileName}}"/>
                <input class="LinkButton clickable" value="" type="submit">
            </form>
        </td>
    {{else}}
        <td></td>
        <td></td>
    {{end}}
        <td>
            <form method="post" action="/gophile">
                <input type="hidden" name="DeleteFile" value="{{.FileName}}"/>
                <input class="DeleteButton clickable" value="" type="submit"/>
            </form>
        </td>
    </tr>
{{end}}
</table>

</body>
</html>`

const IMPRESSUM_HTML = `<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="UTF-8">
    <title>GoPhile Fileserver</title>
    <style type="text/css">
        * {
            font-family: 'Helvetica', 'Arial', sans-serif;
        }

        body {
            background-color: #333;
            height: 100%;
            width: 80%;
            margin: 0 auto;
        }

        header {
            width: 100%;
            height: 50px;
            background-color: #AAA;
        }

        #logo {
            margin-left: 10px;
            height: 50px;
        }

        #main {
            margin-top: 15px;
            background: white;
            padding: 20px;
            border-radius: 15px;
            text-align: center;
        }
    </style>

</head>
<body>
<header>
    <img id="logo"
         src='data:image/png;base64,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'
         alt="gophile">
</header>

<div id="main">
    <h1>Imprint</h1>
    <p>Kevin Lackmann</p>
    <p>Karsten Köhler</p>
    <p>Alexander Nicke</p>
</div>

<footer>
    <p>
        <a href="../">Back to Login</a>
    </p>
</footer>
</body>
</html>`

const INDEX_HTML = `<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="UTF-8">
    <title>GoPhile Fileserver</title>
    <style type="text/css">

        * {
            font-family: 'Helvetica', 'Arial', sans-serif;
        }

        body {
            background-color: #333;
            height: 100%;
            width: 80%;
            margin: 0 auto;
        }

        header {
            width: 100%;
            height: 50px;
            background-color: #AAA;
        }

        #logo {
            margin-left: 10px;
            height: 50px;
        }

        #main {
            margin-top: 15px;
            background: white;
            padding: 20px;
            border-radius: 15px;
            text-align: center;
        }

        fieldset {
            width: 200px;
            display: inline;
            margin: 20px;
            background: #ccc;
            border: none;
            border-radius: 15px;
            padding: 20px;
        }

        legend {
            background: #DDD;
            padding: 5px 15px;
            border: none;
            border-radius: 15px;
        }

        fieldset input, fieldset button {
            width: 100%;
            margin-top: 10px;
        }

    </style>
</head>

<body>

<header>
    <img id="logo"
         src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAO0AAABkCAMAAABdLEetAAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOkAAAGwUExURQAAACMjIyUlJSQkJB4eHiUlJSQkJCUlJSQkJCQkJCUlJSUlJSQkJCUlJSUlJSQkJCYmJiUlJSYmJiUlJSUlJSUlJSUlJSoqKhwcHCUlJSIiIiYmJiUlJSUlJSMjIyQkJCQkJCMjIyUlJSQkJCUlJSUlJSMjIyQkJCQkJCQkJCUlJSQkJCQkJAAAACQkJCUlJSUlJSQkJCUlJSQkJCQkJCYmJgAAACgoKCcnJyUlJSQkJCQkJC4uLigoKCUlJSUlJSUlJScnJyUlJSUlJSUlJSQkJCQkJCUlJSoqKiUlJSQkJCUlJSUlJSQkJCUlJSQkJCUlJSUlJSYmJiUlJT8/PygoKCQkJCMjIyUlJSQkJCMjIyQkJCUlJSYmJiQkJCUlJSUlJSQkJCUlJSYmJiUlJSoqKiUlJSQkJCMjIyUlJScnJyIiIgAAACUlJSQkJCUlJSUlJSQkJCUlJSUlJSUlJSUlJSUlJR8fHyUlJSQkJI+Pj1xcXCIiIiUlJRkZGVBQUDY2NiUlJf///zQ0NOPj48HBwSYmJvz8/KWlpVFRUSwsLLq6uo+Pj1BQUOXl5ZqamjWa31MAAACBdFJOUwAzzN0R7plEd7tmIlWqiHg8UmT+9s7qBgmlFjv9sSSGfTn8rqv1Oqfy+TbKigLsjl7fX9kqSQEmIGAOTAsTgcb3LdWcPtDS+AyW1pVuTTBG79tQwAQZaE++5kCLpC6Fo4fEsnIbErk4HeM0HgPUHOfwVFmXv9y3CHlc+6APugrRXcCa+hQAAAY/SURBVHja7ZplY+Q2EIbNlkzXJJdjxh6VmdsrMzMzM9PVuTL+5Vq0tlib3Vy3u3o/JLvJ2NIzkkejkYPAy8vLy8vLy8vLy8vLy8vLy8usJO/kYBb2n9EVyf+TNm87WVlh28YV+4auyOeYFiKjNlwM2hDDttli0OaEFi4GbUlowYI8twAZRcmC0AYVhFmyKDGZl6f1tJ5Wu6Sfm3SUpw0rAIp83bRuObdE2kRkiWvjhuZrOej6sbG0eZbSVjOxz13joDTTFllMLk7rYpxBymE7FAzZkgcVCz+AuIsp7Hszls2ItuRarRP9cCpoE5Byfc5K1y1Y3YoCOlreLzC3+06yYbQF3902DsegrWKxyylwm8SjVmEn+gUmKlrZL3Vi951gQ2kbgph1M7YmjaalMy0Ydnn02WE656SpGIRsFsYka5NpQ+bQGMJYOSIuNoQW97dmCAVp05U247scki73maZlZFMuHOE5FjUiLZ0EdYVvmlRkENNwTBtMW6C/DgASHCQLN9pGnrnkKbbhlkqrhAVoKIFEg1bzSERxscG0nV3ExZUEXRk70RaykzuQaLghNtUIZJcwXCj+Sbhdxl/uYsP2rqkQRAseSU+LB0iEZW1XJthCN/4UFwrTR/IdRmnGsWG0Ur/iwYbWRFtrqJKYnx2yYqWb+loJ5PfXiuUXt1262zDaWrLi44SWNlRfPQoIlqHVrFMZ33X0NVVNgrQfThcbRisnAxW62k6bqa9mjyXU09aa7skDlWj9gj2WuNow2siWP2tpU93QEn+1pTHbbnT/jIa0Bed46WkoXG0YFdC4N7HRVqZYZJzKuXEnxWUXmT68N+xfxOa5F184eVxrw1qtbDsd3RecleQaRabRq4w7zWJIG+EePv1M9sqrKp9FI5vnd6GbfnJKY2PysROtnNMLgsbCnnFnBoft7dyNfm7+WIxB7C7o917a5pcaG0YbrpcW2mjjKdE+S253YL+il+z3U7TN879Q25hqF9OhbadEu0Jv97mBdqTjG0gLgUnrqxABkZbpsI32519/OntkA2nddrLKfCnU/ZfLbUZ74LN//HNMXjnwwrMZPdXI5re1ta+/VdtMhTZbX4koNTkqHtKOHpc/19bOyJEd2z3UfXgJ2fz+4y/faGwmpm3UqYmLakMMC7lwjhr5AP3lr7W/f5BTTOyyR7oPLx/CPvnue43NxLRcWjaeKkPyUXO0yPLDzzDJSUUSjB+H7SgWv/bRoRMHvzqls5mYtjTlS5ZKXKxNlCthqUaW+8Cnq0dzRTAj82OTNk8Gwzk0GS1+piLddIS5dXBV+1tacYFch9NSYwh6WovN5LSFPlGObBGsVuOyOiTk98qaIgcdtk1ta7WZnFY/Hxtr8SIhxT5hGapSUjoaJp24k5mygJUPac02U6DNNfOxaPV7QWEUh6V2Uv3OxAordh1fOguj4Z6R0hptpkArFboGiWFqi9Zs0kZNgTZNgB6uZHL1PCOFzb7mBrgXXzDtbovNNGjJZIm5gFTyrxwZiqyRIrUGqpORjNb+AXZL3QogiPY8i800aFmJsD8pC2mrTod0wgESPfYC8noClG7haC02U6EdFYBTtD9o2GlO6ngimYDhIRL1mSr/DqHCLTyt2WY6tEGSyT6NwsBZIcjQCVI3A4epsRzP+yPXNhUOXSmt0WZKtNIhrOsRn7FEp3RXXqBtpHygPqI12ASGd1i5VwP0XxR+b2Ex4TsFlb6CqNGQ9twozLU+tW3qFNXzesZpx5PhrCAc/z2WWacNIt0qhWP8mLvmmacNNYtyGa3jFaWZp6W5gJjIk3SjCeaNlmaDadavNfTFi/FrXQPabXt2PvDEkUdndXRRFoZebwGjV0PMsG++/oaS9qL77r3w4CpdCJdncHTzWLErSCvzRe+eWFndW7294xoU0a6/cv+1V++7527hHg8+NpsLkbgLioE1NXk/tp1Q7Hp4ZmMzgP07Yk3lcsmet8ywK4/PdsDC77CWzubvvGeCPbAczJeSJ/WwdywF86atd+lgt9wazKEuuUkJe9tVwVzq/ksVsEevC+ZUd54WWQ8vB3OspQs42NsvC+ZaO7b3rMdu2BrMuy6/mMLeeHOwANoGtnSstyxdEXh5eXl5eXl5eXl5eXn9F/oXh8wNWhBYFTwAAAAASUVORK5CYII='
         alt="gophile">
</header>

<div id="main">
    <h1>Welcome!</h1>
    <p>{{.Error}}</p>
    <!-- forms for login and sign up -->
    <fieldset>
        <legend>Login</legend>
        <form action="/login" method="post">
            <input type="text" placeholder="Username" name="user" required/>
            <input type="password" placeholder="Password" name="password" required style="margin-bottom: 29px"/>

            <input type="submit" value="Login"/>
        </form>
    </fieldset>
    <fieldset>
        <legend>Register</legend>
        <form action="/register" method="post">
            <input type="text" placeholder="Username" name="user" required/>
            <input type="password" placeholder="Password" name="password" required/>
            <input type="password" placeholder="Confirm Password" name="passwordMatch" required>
            <input type="submit" value="Register"/>
        </form>
    </fieldset>
</div>

<footer>
    <p>
        <a href="/impressum">Imprint</a>
    </p>
</footer>

</body>
</html>`
