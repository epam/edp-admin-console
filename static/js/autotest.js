$(function () {

    $('#languageSelection').on('change', function (e) {
        $('.js-form-subsection, .appLangError').hide();
        $($(e.target).data('target')).show();
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

    $('#gitRepoUrl').focusout(function () {
        isGitDataValid();
    });

    $('#repoLogin').focusout(function () {
        let $repoLoginErrEl = $('.repo-login-validation');
        !isFieldValid($(this), REGEX.REPOSITORY_LOGIN)
            ? $repoLoginErrEl.show()
            : $repoLoginErrEl.hide();
    });

    $('#repoPassword').focusout(function () {
        let $repoPasswordErrEl = $('.repo-password-validation');
        !isFieldValid($(this), REGEX.REPOSITORY_PASSWORD)
            ? $repoPasswordErrEl.show()
            : $repoPasswordErrEl.hide();
    });

    $('#nameOfApp').focusout(function () {
        let $appNameErrEl = $('.app-name-validation');
        !isFieldValid($(this), REGEX.APPLICATION_NAME)
            ? $appNameErrEl.show()
            : $appNameErrEl.hide();
    });

    $('#description').focusout(function () {
        let $descriptionErrEl = $('.description-validation');
        !isFieldValid($(this), REGEX.DESCRIPTION)
            ? $descriptionErrEl.show()
            : $descriptionErrEl.hide();
    });

    $('#vcsLogin').focusout(function () {
        let $vcsLoginErrEl = $('.vcs-login-validation');
        !isFieldValid($(this), REGEX.REPOSITORY_LOGIN)
            ? $vcsLoginErrEl.show()
            : $vcsLoginErrEl.hide();
    });

    $('#vcsPassword').focusout(function () {
        let $vcsPasswordErrEl = $('.vcs-password-validation');
        !isFieldValid($(this), REGEX.REPOSITORY_PASSWORD)
            ? $vcsPasswordErrEl.show()
            : $vcsPasswordErrEl.hide();
    });

    $('.create-autotest').click(function (e) {
        e.preventDefault();

        let isGitValid = isGitDataValid();

        let $appNameErrEl = $('.app-name-validation'),
            isAppValid = isFieldValid($('#nameOfApp'), REGEX.APPLICATION_NAME);
        !isAppValid
            ? $appNameErrEl.show()
            : $appNameErrEl.hide();

        let $descriptionErrEl = $('.description-validation'),
            isDescValid = isFieldValid($('#description'), REGEX.DESCRIPTION);
        !isDescValid
            ? $descriptionErrEl.show()
            : $descriptionErrEl.hide();

        let $vcsLoginErrEl = $('.vcs-login-validation'),
            isVcsLoginValid = isFieldValid($('#vcsLogin'), REGEX.REPOSITORY_LOGIN);
        !isVcsLoginValid
            ? $vcsLoginErrEl.show()
            : $vcsLoginErrEl.hide();

        let $vcsPasswordErrEl = $('.vcs-password-validation'),
            isVcsPasswordValid = isFieldValid($('#vcsPassword'), REGEX.REPOSITORY_PASSWORD);
        !isVcsPasswordValid
            ? $vcsPasswordErrEl.show()
            : $vcsPasswordErrEl.hide();

        let isApplicationValid = isApplicationCodeSelected();

        let isVcsBlockValid;
        isVcsBlockValid = $('.vcs-block').length == 0 ? true : isVcsLoginValid && isVcsPasswordValid;

        if (isGitValid && isAppValid && isDescValid && isApplicationValid && isVcsBlockValid) {
            $('#createAutotest').submit();
        }
    });

    $('.tooltip-icon').tooltip();

});

let REGEX = {
    REPOSITORY_LOGIN: /\w/,
    REPOSITORY_PASSWORD: /\w/,
    APPLICATION_NAME: /^[a-z][a-z0-9-]*[a-z0-9]$/,
    GIT_URL: /(?:^git|^ssh|^https?|^git@[-\w.]+):(\/\/)?(.*?)(\.git)(\/?|\#[-\d\w._]+?)$/,
    DESCRIPTION: /^[a-zA-Z0-9]/
};

function isGitDataValid() {
    let $gitUrlErrEl = $('.git-url-validation'),
        isRepoPrivate = $('#isRepoPrivate').is(':checked'),
        $gitRepoUrlEl = $('#gitRepoUrl'),
        $repoLoginEl = $('#repoLogin'),
        $repoPasswordEl = $('#repoPassword'),
        $gitCredsErrEl = $('.git-creds'),
        $repoErrEl = $('.git-repo-error'),
        isValid = false;

    if (isFieldValid($gitRepoUrlEl, REGEX.GIT_URL)) {
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