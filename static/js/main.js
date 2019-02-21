$(function () {

    var CONST = {
        GIT_URL_REGEXP: /(?:git|ssh|https?|git@[-\w.]+):(\/\/)?(.*?)(\.git)(\/?|\#[-\d\w._]+?)$/,
        APP_NAME_REGEXP: /^[a-zA-Z ]+$/,
        REPO_PASS_REGEXP: /\w/,
        REPO_LOGIN_REGEXP: /\w/,
        ROUTE_SITE_REGEXP: /\w/,
        ROUTE_PATH_REGEXP: /(https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9]\.[^\s]{2,})/,
        DB_CAPACITY_REGEXP: /\w/,
        DB_PERSISTENCE_STORAGE_REGEXP: /\w/
    };

    var validationCallbacks = {
        validateGitRepositoryUrl: function () {
            return CONST.GIT_URL_REGEXP.test($(this).val());
        },
        validateNameOfApplication: function () {
            return CONST.APP_NAME_REGEXP.test($(this).val());
        },
        validateRepositoryPassword: function () {
            return CONST.REPO_PASS_REGEXP.test($(this).val());
        },
        validateRepositoryLogin: function () {
            return CONST.REPO_LOGIN_REGEXP.test($(this).val());
        },
        validateRouteSite: function () {
            return CONST.ROUTE_SITE_REGEXP.test($(this).val());
        },
        validateRoutePath: function () {
            return CONST.ROUTE_PATH_REGEXP.test($(this).val());
        },
        validateDbCapacity: function () {
            return CONST.DB_CAPACITY_REGEXP.test($(this).val());
        },
        validateDbPersistentStorage: function () {
            return CONST.DB_PERSISTENCE_STORAGE_REGEXP.test($(this).val());
        }
    };

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

    /*disables other build tools that doesnt related to selected language*/
    /*todo think about better approach*/
    $('#collapseOne .card-body .form__input-wrapper #subFormWrapper .form__radio-btn-wrapper .form__radio-btn').click(function () {
        $(this).parents('.card-body').find('.invalid-feedback').hide()
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

    /*handles main block on button submit*/
    $('.main-info-submit').click(function (e) {
        var $cardBody = $(this).closest('.card-body');
        var isLangChosen = false;
        $.each($cardBody.find('#subFormWrapper input'), function () {
            if ($(this).is(':checked')) {
                isLangChosen = true;
            }
        });

        var $oneBlock = $(this).closest('#collapseOne')
        if (!isLangChosen) {
            e.stopPropagation();
            $oneBlock.find('.card-body .invalid-feedback').show();
            $oneBlock.prev('.card-header').addClass('invalid').removeClass('success');
        } else {
            $oneBlock.find('.card-body .invalid-feedback').hide();
            $oneBlock.prev('.card-header').removeClass('invalid').addClass('success');
        }
    });

    /*handles repo block on button submit*/
    $('.repo-submit').click(function (e) {
        var $cardBody = $(this).closest('.card-body');
        $.each($cardBody.find('div.form-group:not(.hide-element) input'), function () {
            if ($(this).attr('id') === 'gitRepoUrl') {
                _validateInput.bind(this)(validationCallbacks.validateGitRepositoryUrl);
            }
            if ($(this).attr('id') === 'nameOfApp') {
                _validateInput.bind(this)(validationCallbacks.validateNameOfApplication);
            }
            if ($(this).attr('id') === 'repoPassword') {
                _validateInput.bind(this)(validationCallbacks.validateRepositoryPassword);
            }
            if ($(this).attr('id') === 'repoLogin') {
                _validateInput.bind(this)(validationCallbacks.validateRepositoryLogin);
            }
        });

        toggleValidClassOnAccordionTab.bind(this)('#collapseTwo', e);
    });

    /*handles route block on button submit*/
    $('.route-submit').click(function (e) {
        var $cardBody = $(this).closest('.card-body');
        if ($('#needRoute').is(':checked')) {
            $.each($cardBody.find('div.form-group:not(.hide-element) input'), function () {
                if ($(this).attr('id') === 'routeSite') {
                    _validateInput.bind(this)(validationCallbacks.validateRouteSite);
                }
                if ($(this).attr('id') === 'routePath') {
                    _validateInput.bind(this)(validationCallbacks.validateRoutePath);
                }
            });
        }

        toggleValidClassOnAccordionTab.bind(this)('#collapseThree', e);
    });

    /*handles db block on button submit and sends request to add application*/
    $('.db-submit.create-application-submit').click(function (e) {
        var $cardBody = $(this).closest('.card-body');
        if ($('#needDb').is(':checked')) {
            $.each($cardBody.find('div.form-group:not(.hide-element) input'), function () {
                if ($(this).attr('id') === 'dbCapacity') {
                    _validateInput.bind(this)(validationCallbacks.validateDbCapacity);
                }
                if ($(this).attr('id') === 'dbPersitantStorage') {
                    _validateInput.bind(this)(validationCallbacks.validateDbPersistentStorage);
                }
            });
        }

        $('.main-info-submit').trigger('click');
        $('.repo-submit').trigger('click');
        $('.route-submit').trigger('click');
        $('#accordionCreateApplication .card').addClass('collapse')
        toggleValidClassOnAccordionTab.bind(this)('#collapseFour', e);

        if (isFormValid()) {
            var formData = $('#createAppForm').serializeArray();
            var json = buildPayloadToCreateApplication(formData);

            _sendPostRequest(json, function () {
                console.log('Application data is sent to server.');
            }, function () {
                console.log('An error has occurred on server side.');
            });
        }
    });

    /*todo remade button to href*/
    $('div button.create-application.btn-success').click(function () {
        window.location.href = '/admin/application/create'
    });

});

/*service functions*/

function buildPayloadToCreateApplication(formData) {
    var appJson = {
        lang: getValueByName(formData, 'appLang'),
        framework: getValueByName(formData, 'framework'),
        buildTool: getValueByName(formData, 'buildTool'),
        strategy: getValueByName(formData, 'strategy'),
        name: getValueByName(formData, 'nameOfApp')
    };

    if (appJson['strategy'] === 'clone') {
        appJson['git'] = {
            url: getValueByName(formData, 'gitRepoUrl'),
        };

        if (isArrayContainName(formData, 'isRepoPrivate')) {
            appJson['git']['login'] = getValueByName(formData, 'repoLogin');
            appJson['git']['password'] = getValueByName(formData, 'repoPassword');
        }
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
            capacity: getValueByName(formData, 'dbCapacity'),
            storage: getValueByName(formData, 'dbPersitantStorage')
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

function isFormValid() {
    return !($('#accordionCreateApplication').find('.card .card-header.invalid').length > 0);
}

function toggleValidClassOnAccordionTab(tabId, event) {
    var $nonValidInputs = $(this).closest('.card-body').find('div.form-group:not(.hide-element)').find('input.is-invalid');
    if ($nonValidInputs.length > 0) {
        event.stopPropagation();
        $(this).closest(tabId).prev('.card-header').addClass('invalid').removeClass('success');
    } else {
        $(this).closest(tabId).prev('.card-header').removeClass('invalid').addClass('success');
    }
}

function _validateInput(validateCallback) {
    if (!validateCallback.bind(this)()) {
        $(this).next('.invalid-feedback').show();
        $(this).addClass('is-invalid');
    } else {
        $(this).next('.invalid-feedback').hide();
        $(this).removeClass('is-invalid');
    }
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

function _sendPostRequest(data, successCallback, failCallback) {
    $.ajax({
        url: '/api/v1/application/create',
        contentType: "application/json",
        type: "POST",
        data: JSON.stringify(data),
        success: function () {
            successCallback();
        },
        fail: function () {
            failCallback();
        }
    });
}
