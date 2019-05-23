$(function () {

    $('#languageSelection').on('change', function (e) {
        $('.js-form-subsection').hide();
        $($(e.target).data('target')).show();
    });

    $('#languageSelection .form__radio-btn').click(function () {
        $(this).parents('.card-body').find('.appLangError').hide();
        $('.framework').prop('checked', false);

        let $target = $(this).data('target');
        let $frameworks = $(this).parents('.card-body').find('.form-group .form-subsection');

        $.each($frameworks, function () {
            !$(this).hasClass($target.substring(1, $target.length))
                ? $(this).find('select').attr('disabled', true)
                : $(this).find('select').removeAttr('disabled');
        });
    });

    $('#isRepoPrivate').change(function () {
        let $repositoryCredsEl = $('.repo-credentials');
        let $credInputElems = $repositoryCredsEl.find('.repoLogin input, .repoPassword input');
        let $repoErrElems = $('.repo-login-validation, .repo-password-validation');

        if ($(this).is(':checked')) {
            $repositoryCredsEl.removeClass('hide-element');
        } else {
            $repoErrElems.hide();
            $repositoryCredsEl.addClass('hide-element');
            $credInputElems.removeClass('is-invalid');
            $credInputElems.val('');
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

    $('.form__input-wrapper .form-subsection input').click(function () {
        $('.frameworkError').hide();
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

        let isApplicationAndFrameworkValid = isApplicationCodeAndFrameworkSelected();


        if (isGitValid && isAppValid && isDescValid && isApplicationAndFrameworkValid) {
            $('#createAutotest').submit();
        }
    });

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

function isFieldValid(elementToValidate, regex) {
    let check = function (value) {
        return regex.test(value);
    };

    return !(!elementToValidate.val() || !check(elementToValidate.val()));
}

function isApplicationCodeAndFrameworkSelected() {
    let $languageCheckboxElems = $('.language input');
    let $frameworkCheckboxElems = $('.form__input-wrapper .form-subsection input');
    let $frameworkErrEl = $('.frameworkError');
    let $appLanguageErrEl = $('.appLangError');

    if ($languageCheckboxElems.is(':checked')) {
        $appLanguageErrEl.hide();
        if ($frameworkCheckboxElems.is(":checked")) {
            $frameworkErrEl.hide();
            return true;
        } else {
            $frameworkErrEl.show();
            return false;
        }
    } else {
        $appLanguageErrEl.show();
        return false;
    }
}