$(function () {
    !function () {
        var $menuBtn = $('.js-toggle-menu-button');
        var $menu = $('.js-aside-menu');

        $menuBtn.on('click', function () {
            $(this).toggleClass('collapsed');
            $menu.toggleClass('active');
        });
    }();

    !function () {
        var $subFormWrap = $('#subFormWrapper');
        var $subForms = $('.js-form-subsection');

        $subFormWrap.on('change', function (e) {
            var target = $(e.target).data('target');
            $subForms.hide();
            $(target).show();
        });
    }();

    $('[data-toggle="tooltip"]').tooltip();

    $('#addAppForm').submit(function () {
        _sendPostRequest(this, function () {
            console.log('Application data is sent to server.');
        }, function () {
            console.log('An error has occurred on server side.');
        });
    });

    $('div button.btn-success').click(function () {
        window.location.href = '/admin/addingApplication'
    });

});

function _sendPostRequest(that, successCallback, failCallback) {
    $.ajax({
        url: '/api/v1/application/add',
        contentType: "application/json",
        type: "POST",
        data: JSON.stringify(_generateJson(that)),
        success: function () {
            successCallback();
        },
        fail: function () {
            failCallback();
        }
    });
}

function _generateJson(form) {
    var resJson = {
        appLang: '',
        framework: '',
        gitRepoUrl: '',
        strategy: '',
        nameOfApp: '',
        buildTool: '',
        needRoute: '',
        routeSite: '',
        routePath: '',
        database: '',
        dbVersion: '',
        dbCapacity: '',
        dbPersitantStorage: ''
    };
    var targetJson = _formDataToJson(form);

    Object.keys(resJson).forEach(function (key) {
        var value = targetJson[key];
        if (!!value) {
            resJson[key] = value;
        }
    });

    return resJson;
}

function _formDataToJson(form) {
    var json = {};
    $.each($(form).serializeArray(), function () {
        json[this.name] = this.value;
    });

    return json;
}
