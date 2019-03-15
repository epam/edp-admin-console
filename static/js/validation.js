$(function () {

    /*handles repo block on button submit*/
    $('.repo-submit').click(function (event) {
        let isValid = isRepoBlockValid();
        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($repoBlock);
            return;
        }
        blockIsValid($repoBlock);
    });

    /*handles main block on button submit*/
    $('.main-info-submit').click(function (event) {
        let isValid = isMainBlockValid();
        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($mainBlock);
            return;
        }
        blockIsValid($mainBlock);
    });

    /*handles vcs block on button submit*/
    $('.vcs-submit').click(function (event) {
        let isValid = isVcsBlockValid();
        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($vcsBlock);
            return
        }
        blockIsValid($vcsBlock);
    });

    /*handles route block on button submit*/
    $('.route-submit').click(function (event) {
        let isValid = isRouteBlockValid();
        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($routeBlock);
            return
        }
        blockIsValid($routeBlock);
    });

    /*handles db block on button submit and sends request to add application*/
    $('.db-submit').click(function (event) {
        let isValid = true;

        let isDbValid = isDbBlockValid();
        if (!isDbValid) {
            isValid = false;
            event.stopPropagation();
            blockIsNotValid($dbBlock);
        } else {
            blockIsValid($dbBlock);
        }

        let isMainValid = isMainBlockValid();
        if (!isMainValid) {
            isValid = false;
            event.stopPropagation();
            blockIsNotValid($mainBlock);
        } else {
            blockIsValid($mainBlock);
        }

        if (!$('.vcs-block').hasClass('hide-element')) {
            let isVcsValid = isVcsBlockValid();
            if (!isVcsValid) {
                isValid = false;
                event.stopPropagation();
                blockIsNotValid($vcsBlock);
            } else {
                blockIsValid($vcsBlock);
            }
        }

        let isRouteValid = isRouteBlockValid();
        if (!isRouteValid) {
            isValid = false;
            event.stopPropagation();
            blockIsNotValid($routeBlock);
        } else {
            blockIsValid($routeBlock);
        }
        let isRepoValid = isRepoBlockValid();
        if (!isRepoValid) {
            isValid = false;
            event.stopPropagation();
            blockIsNotValid($repoBlock);
        } else {
            blockIsValid($repoBlock);
        }

        if (isValid) {
            createTableWithValue($('#createAppForm').serializeArray());
            $('#confirmationPopup').modal('show');
        }
    });

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
        let $repoUrl = $('#gitRepoUrl');
        let $repoLogin = $('#repoLogin');
        let $repoPassword = $('#repoPassword');
        let $gitCredMsg = $('.git-creds');
        let $gitRepoMsg = $('.git-repo-error');
        $repoUrl.next('.invalid-feedback').hide();
        $gitRepoMsg.hide();
        $repoLogin.next('.invalid-feedback').hide();
        $repoPassword.next('.invalid-feedback').hide();
        $gitCredMsg.hide();
        $repoUrl.removeClass('is-invalid');
        $repoLogin.removeClass('is-invalid');
        $repoPassword.removeClass('is-invalid');
        blockIsValid($repoBlock);
    });

});

let $repoBlock = $('.repo-block');
let $mainBlock = $('.main-block');
let $vcsBlock = $('.vcs-block');
let $routeBlock = $('.route-block');
let $dbBlock = $('.db-block');

let CONST = {
    GIT_URL_REGEXP: /(?:^git|^ssh|^https?|^git@[-\w.]+):(\/\/)?(.*?)(\.git)(\/?|\#[-\d\w._]+?)$/,
    APP_NAME_REGEXP: /^[a-z][a-z0-9-]*[a-z0-9]$/,
    REPO_PASS_REGEXP: /\w/,
    REPO_LOGIN_REGEXP: /\w/,
    ROUTE_SITE_REGEXP: /^[a-z][a-z0-9-]*[a-z0-9]$/,
    ROUTE_PATH_REGEXP: /^\/.*$/,
    DB_CAPACITY_REGEXP: /\w/,
    DB_PERSISTENCE_STORAGE_REGEXP: /\w/
};

let validate = {
    repositoryUrl: function (val) {
        return CONST.GIT_URL_REGEXP.test(val);
    },
    applicationName: function (val) {
        return CONST.APP_NAME_REGEXP.test(val);
    },
    repositoryPassword: function (val) {
        return CONST.REPO_PASS_REGEXP.test(val);
    },
    repositoryLogin: function (val) {
        return CONST.REPO_LOGIN_REGEXP.test(val);
    },
    vcsLogin: function (val) {
        return CONST.REPO_LOGIN_REGEXP.test(val);
    },
    vcsPassword: function (val) {
        return CONST.REPO_PASS_REGEXP.test(val);
    },
    routeSite: function (val) {
        return CONST.ROUTE_SITE_REGEXP.test(val);
    },
    routePath: function (val) {
        return CONST.ROUTE_PATH_REGEXP.test(val);
    },
    dbCapacity: function (val) {
        return CONST.DB_CAPACITY_REGEXP.test(val);
    }
};

function isRepoBlockValid() {
    let $strategy = $('#strategy');
    let isStrategyClone = $strategy.val() === 'clone';
    let $repoUrl = $('#gitRepoUrl');
    let $repoLogin = $('#repoLogin');
    let $repoPassword = $('#repoPassword');
    let $gitCredMsg = $('.git-creds');
    let $repoMsg = $repoUrl.next('.invalid-feedback');
    let $loginMsg = $repoLogin.next('.invalid-feedback');
    let $passMsg = $repoPassword.next('.invalid-feedback');
    let isRepoPrivate = $('#isRepoPrivate').is(':checked');
    let $gitRepoMsg = $('.git-repo-error');
    let isValid = false;
    let refresh = function refreshErrors() {
        $repoMsg.hide();
        $gitRepoMsg.hide();
        $loginMsg.hide();
        $passMsg.hide();
        $gitCredMsg.hide();
        $repoUrl.removeClass('is-invalid');
        $repoLogin.removeClass('is-invalid');
        $repoPassword.removeClass('is-invalid');
    };
    refresh();

    if (isStrategyClone) {
        let isRepoUrlValid = validate.repositoryUrl($repoUrl.val());
        if (!isRepoUrlValid) {
            $repoMsg.show();
            $repoUrl.addClass('is-invalid');
            if (isRepoPrivate) {
                let isLoginValid = validate.repositoryLogin($repoLogin.val());
                if (!isLoginValid) {
                    $loginMsg.show();
                    $repoLogin.addClass('is-invalid');
                    isValid = false;
                }
                let isPasswordValid = validate.repositoryPassword($repoPassword.val());
                if (!isPasswordValid) {
                    $passMsg.show();
                    $repoPassword.addClass('is-invalid');
                    isValid = false;
                }
            }
            isValid = false;
        } else {
            if (isRepoPrivate) {
                let isLoginValid = validate.repositoryLogin($repoLogin.val());
                if (!isLoginValid) {
                    $loginMsg.show();
                    $repoLogin.addClass('is-invalid');
                    isValid = false;
                }
                let isPasswordValid = validate.repositoryPassword($repoPassword.val());
                if (!isPasswordValid) {
                    $passMsg.show();
                    $repoPassword.addClass('is-invalid');
                    isValid = false;
                }

                if (isLoginValid && isPasswordValid) {
                    _sendPostRequest.bind(this)(false, '/api/v1/repository/available', {
                        url: $repoUrl.val(),
                        login: $repoLogin.val(),
                        password: $repoPassword.val()
                    }, function (isAvailable) {
                        if (!isAvailable) {
                            $gitCredMsg.show();
                            $repoUrl.addClass('is-invalid');
                            $repoLogin.addClass('is-invalid');
                            $repoPassword.addClass('is-invalid');
                            isValid = false;
                        } else {
                            isValid = true;
                        }
                    });
                }
            } else {
                _sendPostRequest.bind(this)(false, '/api/v1/repository/available', {url: $repoUrl.val()}, function (isAvailable) {
                    if (!isAvailable) {
                        $gitRepoMsg.show();
                        $repoUrl.addClass('is-invalid');
                        isValid = false;
                    } else {
                        isValid = true;
                    }
                });
            }

        }
    } else {
        isValid = true;
    }
    return isValid;
}

function isMainBlockValid() {
    let isFrameworkChosen;
    let $appName = $mainBlock.find('#nameOfApp');
    let isAppNameValid = validate.applicationName($appName.val());
    let appNameMsg = $appName.parents('.app-name').find('.app-name-validation .regex-error');
    let appDuplicateMsg = $appName.parents('.app-name').find('.app-name-duplicate-validation .duplicate-msg');
    let refresh = function refreshErrors() {
        appNameMsg.hide();
        appDuplicateMsg.hide();
        $('.appLangError').hide();
        $('.frameworkError').hide();
        $appName.removeClass('is-invalid');
    };
    refresh();

    if (!isAppNameValid) {
        appNameMsg.show();
        $appName.addClass('is-invalid');
    }

    let isLangChosen = $mainBlock.find('.language input').is(':checked');
    if (!isLangChosen) {
        $('.appLangError').show();
    } else {
        isFrameworkChosen = $mainBlock.find('.form__input-wrapper .form-subsection input').is(":checked");
        if (!isFrameworkChosen) {
            $('.frameworkError').show();
        }
    }
    return isAppNameValid && isLangChosen && isFrameworkChosen;
}

function isVcsBlockValid() {
    let $vcsLogin = $('#vcsLogin');
    let $vcsPassword = $('#vcsPassword');
    let loginMsg = $vcsLogin.next('.invalid-feedback');
    let passwordMsg = $vcsPassword.next('.invalid-feedback');
    let refresh = function refreshErrors() {
        loginMsg.hide();
        passwordMsg.hide();
        $vcsLogin.removeClass('is-invalid');
        $vcsPassword.removeClass('is-invalid');
    };
    refresh();

    let isVcsLoginValid = validate.vcsLogin($vcsLogin.val());
    if (!isVcsLoginValid) {
        loginMsg.show();
        $vcsLogin.addClass('is-invalid');
    }

    let isVcsPasswordValid = validate.vcsPassword($vcsPassword.val());
    if (!isVcsPasswordValid) {
        passwordMsg.show();
        $vcsPassword.addClass('is-invalid');
    }
    return isVcsLoginValid && isVcsPasswordValid;
}

function isRouteBlockValid() {
    let $routeSite = $('#routeSite');
    let $routePath = $('#routePath');
    let $siteMsg = $routeSite.next('.invalid-feedback');
    let $passwordMsg = $routePath.next('.invalid-feedback');
    let needRoute = $('#needRoute').is(':checked');
    let refresh = function refreshErrors() {
        $siteMsg.hide();
        $passwordMsg.hide();
        $routeSite.removeClass('is-invalid');
        $routePath.removeClass('is-invalid');
    };
    refresh();

    if (needRoute) {
        let isSiteValid = validate.routeSite($routeSite.val());
        if (!isSiteValid) {
            $siteMsg.show();
            $routeSite.addClass('is-invalid');
        }

        let isPathValid = true;
        if ($routePath.val()) {
            isPathValid = validate.routePath($routePath.val());
            if (!isPathValid) {
                $passwordMsg.show();
                $routePath.addClass('is-invalid');
            }
        }

        return isSiteValid && isPathValid;
    }
    return true;
}

function isDbBlockValid() {
    let $dbCapacity = $('#dbCapacity');
    let $capacityMsg = $dbCapacity.parents('.form-group').find('.invalid-feedback');
    let needDb = $('#needDb').is(':checked');
    let refresh = function refreshErrors() {
        $capacityMsg.hide();
        $dbCapacity.removeClass('is-invalid');
    };
    refresh();

    if (needDb) {
        let isCapacityValid = validate.dbCapacity($dbCapacity.val());
        if (!isCapacityValid) {
            $capacityMsg.show();
            $dbCapacity.addClass('is-invalid')
        }
        return isCapacityValid;
    }
    return true;
}

function blockIsNotValid($block) {
    $block.find('.card-header')
        .addClass('invalid')
        .removeClass('success')
        .addClass('error');
}

function blockIsValid($block) {
    $block.find('.card-header')
        .removeClass('invalid')
        .addClass('success')
        .removeClass('error');
}


function getTenantName() {
    var segments = window.location.pathname.split('/');
    if (segments && segments[3]) {
        return segments[3];
    }
    console.error('Couldn\'t get edp name from url.');
}