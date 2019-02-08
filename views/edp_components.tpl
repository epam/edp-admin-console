<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <link rel="stylesheet" href="/static/css/index.css">
    <link rel="stylesheet" href="/static/css/edp-instance.css">
</head>
<body>
<main>
    <header class="edp-header">
        <div class="burger">
            <button class="js-toggle-menu-button" type="button" data-toggle="collapse">
                <span class="opened">&#10006;</span>
                <span class="closed"></span>
            </button>
        </div>
        <div class="logo">
            <img src="/static/img/epam-logo-full-color-pms-colors.png" alt="epam">
        </div>
        <h2>Epam Delivery Platform</h2>
    </header>
    <section class="content d-flex">
        <aside class="p-0 bg-dark active js-aside-menu aside-menu active">
            <nav class="navbar navbar-expand navbar-dark bg-dark flex-column flex-row align-items-start p-0">
                <div class="navbar-collapse">
                    <div class="dropdown">
                      <button class="dropbtn">Select EDP resource</button>
                      <div class="dropdown-content">
                        {{range $k, $v := .EDPComponents}}
                            <a href="{{$v}}">{{$k}}</a>
                        {{end}}
                      </div>
                    </div>
                </div>
            </nav>
        </aside>
        <div class="flex-fill pl-4 pr-4 wrapper">

            <div>
                <table class="table">
                    <thead>
                    <tr>
                        <th>{{ .EDPTenantName }} {{ .EDPVersion }}</th>
                    </tr>
                    </thead>
                    <tbody>
                        {{ range $k, $v := .EDPComponents }}
                            <tr>
                                <td>
                                    <a href="{{ $v }}">{{ $k }}</a>
                                </td>
                            </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        </div>
    </section>
    <footer class="edp-footer">
        <div class="copy">
            <p>© 2018 EPAM Systems, Inc. </p>
            <p> All Rights Reserved.</p>
        </div>
    </footer>
</main>
<script src="../static/js/jquery-3.3.1.js"></script>
<script src="../static/js/popper.js"></script>
<script src="../static/js/bootstrap.js"></script>
<script src="../static/js/main.js"></script>
</body>
</html>