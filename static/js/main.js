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

    !function () {
        $.each($('.route, .dataBase .card-body').find('input, select'), function () {
            if ($(this).is('input')) {
                $(this).attr('readonly', true);
            } else if ($(this).is('select')) {
                $(this).attr('disabled', true);
            }
        });
    }();

    $(document).ready(function () {
        _sendGetRequest('/api/v1/storage-class',
            function (storageClasses) {
                var $select = $('#dbPersistentStorage');

                $.each(storageClasses, function () {
                    $select.append('<option value="' + this.toString() + '">' + this.toString() + '</option>');
                });
            }, function (resp) {
                console.log(resp);
            })
    });

    $('[data-toggle="tooltip"]').tooltip();

    $('#strategy').change(function () {
        if (this.value === 'clone') {
            $(".repo-url, .private-repo").removeClass('hide-element');
            if ($('#isRepoPrivate').is(':checked')) {
                $('.repoLogin, .repoPassword').removeClass('hide-element');
            }
        } else {
            $(".repo-url, .private-repo").addClass('hide-element');
            $('.repoLogin, .repoPassword').addClass('hide-element');
        }
    });

    $('#isRepoPrivate').change(function () {
        if ($(this).is(':checked')) {
            $('.repoLogin, .repoPassword').removeClass('hide-element');
        } else {
            $('.repoLogin, .repoPassword').addClass('hide-element');
            $('.repoLogin, .repoPassword').find('.invalid-feedback').hide();
            $('.repoLogin, .repoPassword').find('input').removeClass('is-invalid');

        }
    });

    /*hides and removes validate classes from optional fields (all route and database)*/
    $('#needRoute, #needDb').change(function () {
        toggleFields.bind(this)();

        $.each($(this).closest('.card-body').find('div.form-group:not(.hide-element) input'), function () {
            $(this).removeClass('is-invalid');
            $(this).next('.invalid-feedback').hide();
        });
    });

    $('.java-build-tools').change(function () {
        if (this.value === 'Maven') {
            $('.multi-module').removeClass('hide-element');
        } else {
            $('.multi-module').addClass('hide-element');
        }
    });

    /*disables other build tools that doesnt related to selected language*/
    /*todo think about better approach*/
    $('#collapseTwo .card-body .form__input-wrapper #subFormWrapper .form__radio-btn-wrapper .form__radio-btn').click(function () {
        $(this).parents('.card-body').find('.appLangError').hide();
        $('.framework').prop('checked', false);
        var $target = $(this).data('target');
        $target = $target.substring(1, $target.length);

        $.each($(this).parents('.card-body').find('.form-group .form-subsection'), function () {
            if (!$(this).hasClass($target)) {
                $(this).find('select').attr('disabled', true);
            } else {
                $(this).find('select').removeAttr('disabled');
            }
        });
    });

    $('#create-application').click(function () {
        $('#createAppForm').submit();
        $('#confirmationPopup').modal('hide');
        $(".window-table-body").remove();
    });

    $("#btn-cross-close, #btn-modal-close").click(function () {
        $(".window-table-body").remove();
    });

});

/*service functions*/

function createTableWithValue($formData) {
    var isVcsEnabled = $('.vcs-block').length !== 0;
    var isNeedRoute = isArrayContainName($formData, 'needRoute');
    var isNeedDb = isArrayContainName($formData, 'needDb');
    var isStrategyClone = getValueByName($formData, 'strategy') === "clone";
    var isRepositoryPrivate = isArrayContainName($formData, 'isRepoPrivate');
    var vcsIntegrationEnabled =  isVcsEnabled ? "&#10004;" : "&#10008;";
    var isAppMultiModule = isArrayContainName($formData, 'isMultiModule') ? "&#10004;" : "&#10008;";
    var table = $("#window-table");

    $('<tbody class="window-table-body">' +
        '<tr><td>EDP Name</td><td>' + getValueByName($formData, 'nameOfApp') + '</td></tr>' +
        '<tr><td>Application language</td><td>' + getValueByName($formData, 'appLang') + '</td></tr>' +
        '<tr><td>Framework</td><td>' + getValueByName($formData, 'framework') + '</td></tr>' +
        '<tr><td>Build tool</td><td>' + getValueByName($formData, 'buildTool') + '</td></tr>' +
        '<tr><td>Integration with VCS is enabled</td><td>' + vcsIntegrationEnabled + '</td></tr>').appendTo(table);

    $('<tr><td>Multi module application</td><td>' + isAppMultiModule + '</td></tr>').appendTo(table);

    $('<tr><td class="font-weight-bold text-center" colspan="2">REPOSITORY</td></tr>' +
        '<tr><td>Strategy</td><td>' + getValueByName($formData, 'strategy') + '</td></tr>').appendTo(table);

    if (isStrategyClone) {
        $('<tr><td>Repository url</td><td>' + getValueByName($formData, 'gitRepoUrl') + '</td></tr>').appendTo(table);

        if (isRepositoryPrivate) {
            $('<tr><td>Login</td><td>' + getValueByName($formData, 'repoLogin') + '</td></tr>').appendTo(table);
        }
    }

    if (isVcsEnabled) {
        $('<tr><td class="font-weight-bold text-center" colspan="2">VCS</td></tr>' +
            '<tr><td>Login</td><td>' + getValueByName($formData, 'vcsLogin') + '</td></tr>').appendTo(table)
    }

    if (isNeedRoute) {
        $('<tr><td class="font-weight-bold text-center" colspan="2">ROUTE</td></tr>' +
            '<tr><td>Route site</td><td>' + getValueByName($formData, 'routeSite') + '</td></tr>' +
            '<tr><td>Route path</td><td>' + getValueByName($formData, 'routePath') + '</td></tr>').appendTo(table)
    }

    if (isNeedDb) {
        $('<tr><td class="font-weight-bold text-center" colspan="2">DATABASE</td></tr>' +
            '<tr><td>Database</td><td>' + getValueByName($formData, 'database') + '</td></tr>' +
            '<tr><td>Version</td><td>' + getValueByName($formData, 'dbVersion') + '</td></tr>' +
            '<tr><td>Capacity</td><td>' + getValueByName($formData, 'dbCapacity') + getValueByName($formData, 'capacityExt') + '</td></tr>' +
            '<tr><td>Persistent storage</td><td>' + getValueByName($formData, 'dbPersistentStorage') + '</td></tr>').appendTo(table)
    }

}

function buildPayloadToCreateApplication(formData) {
    var appJson = {
        lang: getValueByName(formData, 'appLang'),
        framework: getValueByName(formData, 'framework'),
        buildTool: getValueByName(formData, 'buildTool'),
        strategy: getValueByName(formData, 'strategy'),
        name: getValueByName(formData, 'nameOfApp'),
    };

    if(isArrayContainName(formData, 'isMultiModule')) {
        appJson.multiModule = JSON.parse(getValueByName(formData, 'isMultiModule'));
    }

    if (appJson['strategy'] === 'clone') {
        appJson['repository'] = {
            url: getValueByName(formData, 'gitRepoUrl'),
        };

        if (isArrayContainName(formData, 'isRepoPrivate')) {
            appJson['repository']['login'] = getValueByName(formData, 'repoLogin');
            appJson['repository']['password'] = getValueByName(formData, 'repoPassword');
        }

    }

    if (isArrayContainName(formData, 'vcsLogin') && isArrayContainName(formData, 'vcsPassword')) {
        appJson['vcs'] = {
            login: getValueByName(formData, 'vcsLogin'),
            password: getValueByName(formData, 'vcsPassword'),
        };
    }

    if (isArrayContainName(formData, 'needRoute')) {
        appJson['route'] = {
            site: getValueByName(formData, 'routeSite'),
            path: getValueByName(formData, 'routePath')
        };
    }

    if (isArrayContainName(formData, 'needDb')) {
        appJson['database'] = {
            kind: getValueByName(formData, 'database'),
            version: getValueByName(formData, 'dbVersion'),
            capacity: getValueByName(formData, 'dbCapacity') + getValueByName(formData, 'capacityExt'),
            storage: getValueByName(formData, 'dbPersistentStorage')
        };
    }

    return appJson;
}

function getValueByName(array, name) {
    return array.find(x => x.name === name).value
}

function isArrayContainName(array, name) {
    return array.find(x => x.name === name)
}

function toggleFields() {
    var toggleInputs = function (bool) {
        $.each($(this).closest('.card-body').find('input, select'), function () {
            if ($(this).is('input')) {
                $(this).attr('readonly', bool);
            } else if ($(this).is('select')) {
                $(this).attr('disabled', bool);
            }
        })
    }.bind(this);

    if ($(this).is(":checked")) {
        toggleInputs(false);
    } else {
        toggleInputs(true);
    }
}

function _sendPostRequest(async, url, data, successCallback, failCallback) {
    $.ajax({
        url: url,
        contentType: "application/json",
        type: "POST",
        data: JSON.stringify(data),
        async: async,
        success: function (resp) {
            successCallback(resp);
        },
        error: function (resp) {
            failCallback(resp);
        }
    });
}

function _sendGetRequest(url, successCallback, failCallback) {
    $.ajax({
        url: url,
        contentType: "application/json",
        success: function (resp) {
            successCallback(resp);
        },
        error: function (resp) {
            failCallback(resp);
        },
    });
}

function getTenantName() {
    var segments = window.location.pathname.split('/');
    if (segments && segments[3]) {
        return segments[3];
    }
    console.error('Couldn\'t get edp name from url.');
}
