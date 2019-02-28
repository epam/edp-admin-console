$(function () {
    var CONST = {
        GIT_URL_REGEXP: /(?:^git|^ssh|^https?|^git@[-\w.]+):(\/\/)?(.*?)(\.git)(\/?|\#[-\d\w._]+?)$/,
        APP_NAME_REGEXP: /^[a-z]+(-+[a-z0-9]+)*$/,
        REPO_PASS_REGEXP: /\w/,
        REPO_LOGIN_REGEXP: /\w/,
        ROUTE_SITE_REGEXP: /^[a-z][a-z0-9-.]+[a-z]$/,
        ROUTE_PATH_REGEXP: /^\/.*$/,
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
        validateVcsLogin: function () {
            return CONST.REPO_LOGIN_REGEXP.test($(this).val());
        },
        validateVcsPassword: function () {
            return CONST.REPO_PASS_REGEXP.test($(this).val());
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

    $(document).ready(function () {
        _sendGetRequest('/api/v1/' + getTenantName() + '/vcs/',
            function (isVcsEnabled) {
                if (isVcsEnabled) {
                    $('.vcs-block').removeClass('hide-element');
                }
            }, function (resp) {
                console.log(resp);
            });
        _sendGetRequest('/api/v1/storage',
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
    $('#collapseOne .card-body .form__input-wrapper #subFormWrapper .form__radio-btn-wrapper .form__radio-btn').click(function () {
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

    /*handles main block on button submit*/
    $('.main-info-submit').click(function (e) {
        var $cardBody = $(this).closest('.card-body');
        var isLangChosen = false;
        $.each($cardBody.find('#subFormWrapper input'), function () {
            if ($(this).is(':checked')) {
                isLangChosen = true;
            }
        });

        var isFrameworkChosen = false;
        $.each($cardBody.find('.form__input-wrapper .form-subsection input'), function () {
            if ($(this).is(':checked')) {
                isFrameworkChosen = true;
            }
        });

        var $appName = $('#nameOfApp');
        _validateInput.bind($appName)(validationCallbacks.validateNameOfApplication);

        var $mainBlock = $(this).closest('#collapseTwo')
        if (!isLangChosen) {
            e.stopPropagation();
            $mainBlock.find('.card-body .invalid-feedback.appLangError').show();
            $mainBlock.prev('.card-header').addClass('invalid').removeClass('success').addClass('error');
        } else {
            $mainBlock.find('.card-body .invalid-feedback.appLangError').hide();
            $mainBlock.prev('.card-header').removeClass('invalid').addClass('success').removeClass('error');
        }
        if (isLangChosen && !isFrameworkChosen) {
            e.stopPropagation();
            $mainBlock.find('.card-body .invalid-feedback.frameworkError').show();
            $mainBlock.prev('.card-header').addClass('invalid').removeClass('success').addClass('error');
        } else if (isLangChosen && isFrameworkChosen) {
            $mainBlock.find('.card-body .invalid-feedback.frameworkError').hide();
            $mainBlock.prev('.card-header').removeClass('invalid').addClass('success').removeClass('error');
        }

        if ($appName.hasClass('is-invalid')) {
            $mainBlock.prev('.card-header').addClass('invalid').removeClass('success').addClass('error');
        } else {
            $mainBlock.prev('.card-header').removeClass('invalid').addClass('success').removeClass('error');
        }

        toggleValidClassOnAccordionTab.bind(this)('#collapseTwo', e);
    });

    /*handles repo block on button submit*/
    $('.repo-submit').click(function (e) {
        $('.repo-validation .invalid-feedback.git-repo').hide();
        var $cardBody = $(this).closest('.card-body');

        if ($('#strategy').val() === 'clone') {
            e.stopPropagation();
            if ($('#gitRepoUrl').val()) {
                var json = {};
                if ($('#isRepoPrivate').is(':checked')) {
                    json.url = $('#gitRepoUrl').val();
                    json.login = $('#repoLogin').val();
                    json.password = $('#repoPassword').val();

                    if (!json.login || !json.password) {
                        var $login = $('#repoLogin');
                        var $password = $('#repoPassword');
                        if (!json.login) {
                            $login.addClass('is-invalid');
                            $login.next('.invalid-feedback').show();
                            $(this).closest('#collapseOne').prev('.card-header').addClass('invalid').removeClass('success').addClass('error');
                        } else {
                            $login.removeClass('is-invalid');
                            $login.next('.invalid-feedback').hide();
                        }

                        if (!json.password) {
                            $password.addClass('is-invalid');
                            $password.next('.invalid-feedback').show();
                            $(this).closest('#collapseOne').prev('.card-header').addClass('invalid').removeClass('success').addClass('error');
                        } else {
                            $password.removeClass('is-invalid');
                            $password.next('.invalid-feedback').hide();
                        }
                        return;
                    }

                    _sendPostRequest.bind(this)('/api/v1/repository', json, function (isAvailable) {
                        if (isAvailable) {
                            $('#gitRepoUrl').removeClass('is-invalid');
                            $('#repoLogin').removeClass('is-invalid');
                            $('#repoPassword').removeClass('is-invalid');

                            $(this).closest('#collapseOne').prev('.card-header').removeClass('invalid').addClass('success').removeClass('error');
                            $('.repo-validation .invalid-feedback.git-creds').hide();

                            $('#headingTwo').trigger('click');
                            $('#headingOne').trigger('click');
                        } else {
                            $('#gitRepoUrl').addClass('is-invalid');
                            $('#repoLogin').addClass('is-invalid');
                            $('#repoPassword').addClass('is-invalid');

                            $(this).closest('#collapseOne').prev('.card-header').addClass('invalid').removeClass('success').addClass('error');
                            $('.repo-validation .invalid-feedback.git-creds').show();
                        }
                    }.bind(this), function () {
                        console.error('An error has occurred while checking existing repository.');
                    });
                } else {
                    json.url = $('#gitRepoUrl').val();

                    _sendPostRequest.bind(this)('/api/v1/repository', json, function (isAvailable) {
                        if (isAvailable) {
                            $('#gitRepoUrl').removeClass('is-invalid');
                            $('#repoLogin').removeClass('is-invalid');
                            $('#repoPassword').removeClass('is-invalid');

                            $(this).closest('#collapseOne').prev('.card-header').removeClass('invalid').addClass('success').removeClass('error');
                            $('.repo-validation .invalid-feedback.git-repo').hide();

                            $('#headingTwo').trigger('click');
                            $('#headingOne').trigger('click');
                        } else {
                            $('#gitRepoUrl').addClass('is-invalid');

                            if (json.login) {
                                $('#repoLogin').addClass('is-invalid');
                                $('#repoPassword').addClass('is-invalid');
                            }
                            $('.repo-validation .invalid-feedback.git-repo').show();
                            $(this).closest('#collapseOne').prev('.card-header').addClass('invalid').removeClass('success').addClass('error');
                        }
                    }.bind(this), function () {
                        console.error('An error has occurred while checking existing repository.');
                    });
                }
            }
        }

        $.each($cardBody.find('div.form-group:not(.hide-element) input'), function () {
            if ($(this).attr('id') === 'gitRepoUrl') {
                _validateInput.bind(this)(validationCallbacks.validateGitRepositoryUrl);
            }
            if ($(this).attr('id') === 'repoPassword') {
                _validateInput.bind(this)(validationCallbacks.validateRepositoryPassword);
            }
            if ($(this).attr('id') === 'repoLogin') {
                _validateInput.bind(this)(validationCallbacks.validateRepositoryLogin);
            }
        });

        toggleValidClassOnAccordionTab.bind(this)('#collapseOne', e);
    });

    /*handles vcs block on button submit*/
    $('.vcs-submit').click(function (e) {
        var $cardBody = $(this).closest('.card-body');
        $.each($cardBody.find('div.form-group:not(.hide-element) input'), function () {
            if ($(this).attr('id') === 'vcsLogin') {
                _validateInput.bind(this)(validationCallbacks.validateVcsLogin);
            }
            if ($(this).attr('id') === 'vcsPassword') {
                _validateInput.bind(this)(validationCallbacks.validateVcsPassword);
            }
        });

        toggleValidClassOnAccordionTab.bind(this)('#collapseThree', e);
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

        toggleValidClassOnAccordionTab.bind(this)('#collapseFour', e);
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
        $('.vcs-submit').trigger('click');
        $('.repo-submit').trigger('click');
        $('.route-submit').trigger('click');
        toggleValidClassOnAccordionTab.bind(this)('#collapseFive', e);

        if (isFormValid()) {
            var $confTable = $('#confirmationPopup').find('.modal-body .confirmation-data');
            setConfirmationData($('#createAppForm').serializeArray(), $confTable);

            $('#confirmationPopup').modal('show');
        }
    });

    $('.create-application').click(function () {
        var json = buildPayloadToCreateApplication($('#createAppForm').serializeArray());
        $('#confirmationPopup').modal('hide');
        /*_sendPostRequest('/api/v1/' + getTenantName() + '/application/create', json, function () {
            $('#successPopup').modal('show');
        }, function () {
            $('#errorPopup').modal('show');
        });*/
    });

});

/*service functions*/

function setConfirmationData(json, $confTable) {
    $confTable.empty();
    $confTable.append($('<tr>')).append('Application language:' + getValueByName(json, 'appLang'));
    $confTable.append($('<tr>')).append('Framework:' + getValueByName(json, 'framework'));
    $confTable.append($('<tr>')).append('Build tool:' + getValueByName(json, 'buildTool'));

    if(isArrayContainName(json, 'isMultiModule')) {
        $confTable.append($('<tr>')).append('Is multi module:' + true);
    }

    $confTable.append($('<tr>')).append('Strategy:' + getValueByName(json, 'strategy'));
    $confTable.append($('<tr>')).append('Application name' + getValueByName(json, 'nameOfApp'));

    if (json['strategy'] === 'clone') {
        $confTable.append($('<tr>')).append('Git url:' + getValueByName(json, 'gitRepoUrl'));
        if (isArrayContainName(json, 'isRepoPrivate')) {
            $confTable.append($('<tr>')).append('Git login:' + getValueByName(json, 'repoLogin'));
        }
    }

    if (isArrayContainName(json, 'vcsLogin') && isArrayContainName(json, 'vcsPassword')) {
        $confTable.append($('<tr>')).append('Vcs login:' + getValueByName(json, 'vcsLogin'));
    }

    if (isArrayContainName(json, 'needRoute')) {
        $confTable.append($('<tr>')).append('Route site:' + getValueByName(json, 'routeSite'));
        $confTable.append($('<tr>')).append('Route path:' + getValueByName(json, 'routePath'));
    }

    if (isArrayContainName(json, 'needDb')) {
        $confTable.append($('<tr>')).append('Database:' + getValueByName(json, 'database'));
        $confTable.append($('<tr>')).append('Database version:' + getValueByName(json, 'dbVersion'));
        $confTable.append($('<tr>')).append('Database capacity:' + getValueByName(json, 'dbCapacity') + getValueByName(json, 'capacityExt'));
        $confTable.append($('<tr>')).append('Database persistent storage:' + getValueByName(json, 'dbPersistentStorage'));
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

function isFormValid() {
    return !($('#accordionCreateApplication').find('.card .card-header.invalid').length > 0);
}

function toggleValidClassOnAccordionTab(tabId, event) {
    var $nonValidInputs = $(this).closest('.card-body').find('div.form-group:not(.hide-element)').find('input.is-invalid');
    if ($nonValidInputs.length > 0) {
        event.stopPropagation();
        $(this).closest(tabId).prev('.card-header').addClass('invalid').removeClass('success').addClass('error');
    } else {
        $(this).closest(tabId).prev('.card-header').removeClass('invalid').addClass('success').removeClass('error');
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

function _sendPostRequest(url, data, successCallback, failCallback) {
    $.ajax({
        url: url,
        contentType: "application/json",
        type: "POST",
        data: JSON.stringify(data),
        success: function (resp) {
            successCallback(resp);
        },
        fail: function () {
            failCallback();
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
        fail: function (resp) {
            failCallback(resp);
        }
    });
}

function getTenantName() {
    var segments = window.location.pathname.split('/');
    if (segments && segments[2]) {
        return segments[2];
    }
    console.error('Couldn\'t get edp name from url.');
}
