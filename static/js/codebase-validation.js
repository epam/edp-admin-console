$(function () {
    let REGEX = {
        CAPACITY: /\w/,
        SERVICE_PATH: /^\/.*$/,
        SERVICE_NAME: /^[a-z][a-z0-9-]*[a-z0-9]$/,
        VCS_LOGIN: /\w/,
        VCS_PASSWORD: /\w/,
        DESCRIPTION: /^[a-zA-Z0-9]/,
        CODEBASE_NAME: /^[a-z][a-z0-9-]*[a-z0-9]$/,
        REPO_LOGIN: /\w/,
        REPO_PASSWORD: /\w/,
        REPO_URL: /(?:^git|^ssh|^https?|^git@[-\w.]+):(\/\/)?(.*?)(\.git)(\/?|\#[-\d\w._]+?)$/
    };

    $('.tooltip-icon').add('[data-toggle="tooltip"]').tooltip();

    !function () {
        $('.form-group .js-form-subsection select').attr('disabled', true);

        $('.multi-module').addClass('hide-element');
        $('#multiModule').attr("disabled", true);
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

    !function () {
        if ($('.db-block').length !== 0) {
            _sendGetRequest(true, '/api/v1/storage-class',
                function (storageClasses) {
                    var $select = $('#dbPersistentStorage');

                    $.each(storageClasses, function () {
                        $select.append('<option value="' + this.toString() + '">' + this.toString() + '</option>');
                    });
                }, function (resp) {
                    console.log(resp);
                })
        }
    }();

    $('#languageSelection').on('change', function (e) {
        $('.js-form-subsection, .appLangError').hide();
        let langDivEl = $($(e.target).data('target'));
        langDivEl.find('select').attr('disabled', false);
        langDivEl.show();
        $('.js-form-subsection input[name="framework"]').prop('checked', false);

        $('.multi-module').addClass('hide-element');
        $('#multiModule').attr("disabled", true);

        let codebaseVal = $('.card.main-block').data('codebase-type');
        if (codebaseVal === 'application' || codebaseVal === 'library') {
            $('.java-build-tools').val('Gradle');
        } else {
            $('.java-build-tools').val('Maven');
        }

        $('.test-report-framework').val('allure');
    });

    $('#isRepoPrivate').change(function () {
        let $login = $('.repoLogin'),
            $pass = $('.repoPassword');
        if ($(this).is(':checked')) {
            $login.add($pass).removeClass('hide-element');
        } else {
            $login.add($pass).addClass('hide-element');
            $login.add($pass).find('.invalid-feedback').hide();
            $login.add($pass).find('input').removeClass('is-invalid');
        }
    });

    $('#strategy').change(function () {
        let $login = $('.repoLogin'), $pass = $('.repoPassword'),
            $url = $('.repo-url'), $privateRep = $('.private-repo');

        if (this.value === 'clone') {
            $url.add($privateRep).removeClass('hide-element');
            if ($('#isRepoPrivate').is(':checked')) {
                $login.add($pass).removeClass('hide-element');
            }
        } else {
            $url.add($privateRep).addClass('hide-element');
            $login.add($pass).addClass('hide-element');
        }
    });

    $('#btn-modal-continue').click(function () {
        $('form.edp-form').submit();
        $('#confirmationPopup').modal('hide');
        $(".window-table-body").remove();
    });

    $("#btn-cross-close, #btn-modal-close").click(function () {
        $(".window-table-body").remove();
    });

    $('#needRoute').change(function () {
        let $exposeServiceBlockEl = $('.route-block'),
            $inputsEl = $exposeServiceBlockEl.find('input');

        if ($(this).is(":checked")) {
            $inputsEl.attr('readonly', false);
        } else {
            $inputsEl.attr('readonly', true);
        }

        $inputsEl.removeClass('is-invalid').next('.invalid-feedback').hide();
    });

    $('#needDb').change(function () {
        let $dbBlockEl = $('.db-block'),
            $inputsEl = $dbBlockEl.find('input'),
            $selectsEl = $dbBlockEl.find('select');

        if ($(this).is(":checked")) {
            $inputsEl.attr('readonly', false);
            $selectsEl.attr('disabled', false);
        } else {
            $inputsEl.attr('readonly', true);
            $selectsEl.attr('disabled', true);
        }

        $('.capacity-error.invalid-feedback').hide();
        $inputsEl.removeClass('is-invalid');
    });

    $('.codebase-info-button').click(function (event) {
        validateCodebaseInfo(event);
    });

    $('.application-submit,.autotest-submit,.library-submit').click(function (event) {
        validateMainInfo(event);
    });

    $('.vcs-submit,.create-library,.create-autotest').click(function (event) {
        if ($(this).hasClass('create-autotest') || $(this).hasClass('create-library')) {
            event.preventDefault();

            let canCreateAutotest = validateCodebaseInfo(event) & validateMainInfo(event) & validateVCSInfo(event);
            if (canCreateAutotest) {
                createConfirmTable($(this).hasClass('create-autotest') ? '#createAutotest' : '#createLibrary');
                $('#confirmationPopup').modal('show');
            }
        } else {
            validateVCSInfo(event);
        }
    });

    $('.route-submit').click(function (event) {
        validateRouteInfo(event);
    });

    $('.db-submit').click(function (event) {
        let canCreateApplication = validateCodebaseInfo(event) &
            validateMainInfo(event) & validateVCSInfo(event) &
            validateRouteInfo(event) & validateDbInfo(event);
        if (canCreateApplication) {
            createConfirmTable('#createAppForm');
            $('#confirmationPopup').modal('show');
        }
    });

    $('.java-build-tools').change(function () {
        if (this.value === 'Maven') {
            $('#multiModule').attr("disabled", false);
            $('.multi-module').removeClass('hide-element');
        } else {
            $('.multi-module').addClass('hide-element');
            $('#multiModule').attr("disabled", true);
        }
    });

    function validateCodebaseInfo(event) {
        let $codebaseBlockEl = $('.codebase-block');

        resetErrors($codebaseBlockEl);

        let isValid = isCodebaseInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($codebaseBlockEl);
            return isValid;
        }
        blockIsValid($codebaseBlockEl);

        return isValid;
    }

    function validateMainInfo(event) {
        let $mainBlockEl = $('.main-block');

        resetErrors($mainBlockEl);

        let isValid = isMainInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($mainBlockEl);
            return isValid;
        }
        blockIsValid($mainBlockEl);

        return isValid;
    }

    function validateVCSInfo(event) {
        let $vcsBlockEl = $('.vcs-block');

        resetErrors($vcsBlockEl);

        let isValid = $vcsBlockEl.length === 0 ? true : isVCSValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($vcsBlockEl);
            return isValid;
        }
        blockIsValid($vcsBlockEl);

        return isValid;
    }

    function validateRouteInfo(event) {
        let $exposeServiceBlockEl = $('.route-block');

        resetErrors($exposeServiceBlockEl);

        let isValid = isExposingServiceInfoValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($exposeServiceBlockEl);
            return isValid;
        }
        blockIsValid($exposeServiceBlockEl);

        return isValid;
    }

    function validateDbInfo(event) {
        let $dbBlockEl = $('.db-block');

        resetErrors($dbBlockEl);

        let isValid = isDatabaseValid();

        if (!isValid) {
            event.stopPropagation();
            blockIsNotValid($dbBlockEl);
            return isValid;
        }
        blockIsValid($dbBlockEl);

        return isValid;
    }

    function resetErrors($el) {
        $el.find('input.is-invalid').removeClass('is-invalid');
        $el.find('.invalid-feedback').hide();
    }

    function isCodebaseInfoValid() {
        let isValid = true;
        let $codebaseBlockEl = $('.codebase-block'),
            $strategyEl = $codebaseBlockEl.find('#strategy');

        if ($strategyEl.length === 0 || $strategyEl.val() === 'clone') {
            let $repoUrl = $('#gitRepoUrl'),
                isRepoUrlValid = isFieldValid($repoUrl, REGEX.REPO_URL),
                $repoMsg = $repoUrl.next('.invalid-feedback');

            if (isRepoUrlValid) {
                let $gitRepoMsg = $('.git-repo-error');

                let creds = {
                    url: $repoUrl.val()
                };

                let isRepoPrivate = $('#isRepoPrivate').is(':checked'),
                    $repoLogin = $('#repoLogin'),
                    $repoPassword = $('#repoPassword');
                if (isRepoPrivate) {
                    let isLoginValid = isFieldValid($repoLogin, REGEX.REPO_LOGIN);
                    if (!isLoginValid) {
                        $repoLogin.next('.invalid-feedback').show();
                        $repoLogin.addClass('is-invalid');
                        isValid = false;
                    }
                    let isPasswordValid = isFieldValid($repoPassword, REGEX.REPO_PASSWORD);
                    if (!isPasswordValid) {
                        $repoPassword.next('.invalid-feedback').show();
                        $repoPassword.addClass('is-invalid');
                        isValid = false;
                    }

                    if (isLoginValid && isPasswordValid) {
                        creds.login = $repoLogin.val();
                        creds.password = $repoPassword.val();
                    }
                }

                _sendPostRequest.bind(this)(false, '/api/v1/repository/available', creds, function (isAvailable) {
                    if (isRepoPrivate) {
                        if (isAvailable) {
                            isValid = true;
                        } else {
                            $('.git-creds').show();
                            $repoUrl.addClass('is-invalid');
                            $repoLogin.addClass('is-invalid');
                            $repoPassword.addClass('is-invalid');
                            isValid = false;
                        }
                    } else {
                        if (isAvailable) {
                            isValid = true;
                        } else {
                            $gitRepoMsg.show();
                            $repoUrl.addClass('is-invalid');
                            isValid = false;
                        }
                    }
                });

            } else {
                isValid = false;
                $repoMsg.show();
                $repoUrl.addClass('is-invalid');
            }
        }
        return isValid;
    }

    function isMainInfoValid() {
        let $codebaseEl = $('.main-block'),
            $codebaseInputEl = $('.codebase-name'),
            isCodebaseNameValid = isFieldValid($codebaseInputEl, REGEX.CODEBASE_NAME);

        if (!isCodebaseNameValid) {
            $('.codebase-name-validation.regex-error').show();
            $codebaseInputEl.addClass('is-invalid');
        }

        let $descriptionInputEl = $('#description'),
            $descriptionErrEl = $('.description-validation.regex-error'),
            isDescriptionValid = $descriptionInputEl.length === 0 ? true : isFieldValid($descriptionInputEl, REGEX.DESCRIPTION);

        if (!isDescriptionValid) {
            $descriptionErrEl.show();
            $descriptionInputEl.addClass('is-invalid');
        }

        let isLanguageChosen = $codebaseEl.find('.language input').is(':checked'),
            isFrameworkChosen;
        if (isLanguageChosen) {
            let $frameworksEl = $codebaseEl.find('.form__input-wrapper .form-subsection input');
            isFrameworkChosen = $frameworksEl.length === 0 ? true : $frameworksEl.is(":checked");
            if (!isFrameworkChosen) {
                $('.frameworkError').show();
            }
        } else {
            $('.appLangError').show();
        }

        return isCodebaseNameValid && isDescriptionValid && isLanguageChosen && isFrameworkChosen;
    }

    function isVCSValid() {
        let $vcsLoginInputEl = $('#vcsLogin'),
            isVcsLoginValid = isFieldValid($vcsLoginInputEl, REGEX.VCS_LOGIN);

        if (!isVcsLoginValid) {
            $('.invalid-feedback.vcs-login-validation').show();
            $vcsLoginInputEl.addClass('is-invalid');
        }

        let $vcsPasswordInputEl = $('#vcsPassword'),
            isVcsPasswordValid = isFieldValid($vcsPasswordInputEl, REGEX.VCS_PASSWORD);

        if (!isVcsPasswordValid) {
            $('.invalid-feedback.vcs-password-validation').show();
            $vcsPasswordInputEl.addClass('is-invalid');
        }

        return isVcsLoginValid && isVcsPasswordValid;
    }

    function isExposingServiceInfoValid() {
        let needRoute = $('#needRoute').is(':checked');

        if (needRoute) {
            let $serviceNameInputEl = $('#routeSite'),
                isServiceNameValid = isFieldValid($serviceNameInputEl, REGEX.SERVICE_NAME);

            if (!isServiceNameValid) {
                $('.route-site.invalid-feedback').show();
                $serviceNameInputEl.addClass('is-invalid');
            }

            let $servicePathInputEl = $('#routePath'),
                isServicePathValid = isFieldValid($servicePathInputEl, REGEX.SERVICE_PATH);

            if (!isServicePathValid) {
                $('.route-path.invalid-feedback').show();
                $servicePathInputEl.addClass('is-invalid');
            }

            return isServiceNameValid && isServicePathValid;
        }

        return true;
    }

    function isDatabaseValid() {
        let needDb = $('#needDb').is(':checked');

        if (needDb) {
            let $capacityInputEl = $('#dbCapacity'),
                isCapacityValid = isFieldValid($capacityInputEl, REGEX.CAPACITY);

            if (!isCapacityValid) {
                $('.capacity-error.invalid-feedback').show();
                $capacityInputEl.addClass('is-invalid');
            }

            return isCapacityValid;
        }

        return true;
    }
});