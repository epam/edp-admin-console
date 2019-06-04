$(function () {

    let REGEX = {
        APPLICATION_NAME: /^[a-z][a-z0-9-]*[a-z0-9]$/,
        GIT_URL: /(?:^git|^ssh|^https?|^git@[-\w.]+):(\/\/)?(.*?)(\.git)(\/?|\#[-\d\w._]+?)$/,
        REPOSITORY_LOGIN: /\w/,
        REPOSITORY_PASSWORD: /\w/
    };

    function validateInputField($el, $errEl, regex) {
        let isAppValid = isFieldValid($el, regex);
        $errEl.toggle(!isAppValid);
        return isAppValid;
    }

    function isGitDataValid() {
        let $gitUrlErrEl = $('.git-url-validation'),
            isRepoPrivate = $('#isRepoPrivate').is(':checked'),
            $gitRepoUrlEl = $('#gitRepoUrl'),
            $repoLoginEl = $('#repoLogin'),
            $repoPasswordEl = $('#repoPassword'),
            $gitCredsErrEl = $('.git-creds'),
            $repoErrEl = $('.git-repo-error'),
            isValid = false;

        if (validateInputField($(this), $gitUrlErrEl, REGEX.GIT_URL)) {
            $gitUrlErrEl.hide();

            if (isRepoPrivate) {
                _sendPostRequest.bind(this)(false, '/api/v1/repository/available', {
                    url: $gitRepoUrlEl.val(),
                    login: $repoLoginEl.val(),
                    password: $repoPasswordEl.val()
                }, function (isAvailable) {
                    isValid = !!isAvailable;
                });

                if (isValid) {
                    $gitCredsErrEl.add($repoErrEl).hide();
                    $gitRepoUrlEl.add($repoLoginEl).add($repoPasswordEl).removeClass('is-invalid');
                } else {
                    $gitCredsErrEl.show();
                    $gitRepoUrlEl.add($repoLoginEl).add($repoPasswordEl).addClass('is-invalid');
                }
            } else {
                _sendPostRequest.bind(this)(false, '/api/v1/repository/available', {url: $gitRepoUrlEl.val()},
                    function (isAvailable) {
                        isValid = !!isAvailable;
                    });

                if (isValid) {
                    $repoErrEl.hide();
                    $gitRepoUrlEl.removeClass('is-invalid');
                } else {
                    $gitCredsErrEl.hide();
                    $repoErrEl.show();
                    $gitRepoUrlEl.addClass('is-invalid');
                }
            }
        } else {
            $gitUrlErrEl.show();
        }
        return isValid;
    }

    function isApplicationCodeSelected() {
        let $languageCheckboxElems = $('.language input');
        let $appLanguageErrEl = $('.appLangError');

        if ($languageCheckboxElems.is(':checked')) {
            $appLanguageErrEl.hide();
            return true;
        }
        $appLanguageErrEl.show();
        return false;
    }

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

    $('#isRepoPrivate').change(function () {
        let $login = $('.repoLogin'), $pass = $('.repoPassword');
        if ($(this).is(':checked')) {
            $login.add($pass).removeClass('hide-element');
        } else {
            $login.add($pass).addClass('hide-element');
            $login.add($pass).find('.invalid-feedback').hide();
            $login.add($pass).find('input').removeClass('is-invalid');

        }
    });

    $('#languageSelection').on('change', function (e) {
        $('.js-form-subsection, .appLangError').hide();
        $($(e.target).data('target')).show();
    });

    $('#nameOfApp').focusout(function () {
        validateInputField($(this), $('.app-name-validation'), REGEX.APPLICATION_NAME);
    });

    $('#gitRepoUrl').focusout(function () {
        let isValid = isGitDataValid.bind(this)();
        $('.repo-login-validation, .repo-password-validation').toggle(!isValid);
    });

    $('#repoLogin').focusout(function () {
        validateInputField($(this), $('.repo-login-validation'), REGEX.REPOSITORY_LOGIN);
    });

    $('#repoPassword').focusout(function () {
        validateInputField($(this), $('.repo-password-validation'), REGEX.REPOSITORY_PASSWORD);
    });

    $('#vcsLogin').focusout(function () {
        validateInputField($(this), $('.vcs-login-validation'), REGEX.REPOSITORY_LOGIN);
    });

    $('#vcsPassword').focusout(function () {
        validateInputField($(this), $('.vcs-password-validation'), REGEX.REPOSITORY_PASSWORD);
    });

    $('.codebase-info-button').click(function (e) {
        let isStrategyClone = $('#strategy').val() === 'clone',
            isRepoPrivate = $('#isRepoPrivate').is(':checked'),
            isGitValid = true,
            isRepoLoginValid = true,
            isRepoPasswordValid = true,
            $codebaseBlockEl = $('.codebase-block');

        if (isStrategyClone) {
            isGitValid = isGitDataValid.bind($('#gitRepoUrl'))();
            if (isRepoPrivate) {
                isRepoLoginValid = validateInputField($('#repoLogin'), $('.repo-login-validation'), REGEX.REPOSITORY_LOGIN);
                isRepoPasswordValid = validateInputField($('#repoPassword'), $('.repo-password-validation'), REGEX.REPOSITORY_PASSWORD);
            }
        }

        if (!(isGitValid && isRepoLoginValid && isRepoPasswordValid)) {
            e.stopPropagation();
            blockIsNotValid($codebaseBlockEl);
            return;
        }
        blockIsValid($codebaseBlockEl);
    });

    $('.library-submit').click(function (e) {
        let isApplicationFieldValid = validateInputField($('#nameOfApp'), $('.app-name-validation'), REGEX.APPLICATION_NAME),
            isLanguageSelected = isApplicationCodeSelected(),
            $libraryBlockEl = $('.library-block');

        if (!(isApplicationFieldValid && isLanguageSelected)) {
            e.stopPropagation();
            blockIsNotValid($libraryBlockEl);
            return;
        }
        blockIsValid($libraryBlockEl);
    });

    $('.create-library').click(function (e) {
        e.preventDefault();

        let isStrategyClone = $('#strategy').val() === 'clone',
            isRepoPrivate = $('#isRepoPrivate').is(':checked'),
            isGitValid = true,
            isRepoLoginValid = true,
            isRepoPasswordValid = true,
            isApplicationFieldValid = validateInputField($('#nameOfApp'), $('.app-name-validation'), REGEX.APPLICATION_NAME),
            isLanguageSelected = isApplicationCodeSelected();

        if (isStrategyClone) {
            isGitValid = isGitDataValid.bind($('#gitRepoUrl'))();
            if (isRepoPrivate) {
                isRepoLoginValid = validateInputField($('#repoLogin'), $('.repo-login-validation'), REGEX.REPOSITORY_LOGIN);
                isRepoPasswordValid = validateInputField($('#repoPassword'), $('.repo-password-validation'), REGEX.REPOSITORY_PASSWORD);
            }
        }

        let isVcsLoginValid = validateInputField($('#vcsLogin'), $('.vcs-login-validation'), REGEX.REPOSITORY_LOGIN),
            isVcsPasswordValid = validateInputField($('#vcsPassword'), $('.vcs-password-validation'), REGEX.REPOSITORY_PASSWORD),
            isVcsBlockValid = $('.vcs-block').length === 0 ? true : isVcsLoginValid && isVcsPasswordValid;


        let $vcsBlockEl = $('.vcs-block'),
            $libraryBlockEl = $('.library-block'),
            $codebaseBlockEl = $('.codebase-block');

        !(isGitValid && isRepoLoginValid && isRepoPasswordValid)
            ? blockIsNotValid($codebaseBlockEl)
            : blockIsValid($codebaseBlockEl);
        !(isApplicationFieldValid && isLanguageSelected)
            ? blockIsNotValid($libraryBlockEl)
            : blockIsValid($libraryBlockEl);
        !isVcsBlockValid
            ? blockIsNotValid($vcsBlockEl)
            : blockIsValid($vcsBlockEl);

        if (isApplicationFieldValid && isLanguageSelected
            && (isGitValid && isRepoLoginValid && isRepoPasswordValid)
            && isVcsBlockValid) {
            $('#createLibrary').submit();
        }
    });

});