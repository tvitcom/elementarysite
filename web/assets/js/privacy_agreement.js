	var agreePrivacyPolicyOnce;
	var clientParams;

	document.onreadystatechange = function () {
	  if (document.readyState === 'complete') {

	  	agreePrivacyPolicyOnce = function() {
			when = 31536000; // Expected for year
		    document.cookie = "privacy_signer_time="+Date.now()+"; Path=/; secure; SameSite=strict; max-age="+when;
		    document.cookie = "privacy_signer_ua="+btoa(navigator.userAgent)+"; Path=/; secure; SameSite=strict; max-age="+when;
		    document.cookie = "privacy_signer_screen="+screen.availWidth+"x"+screen.availHeight+"; Path=/; SameSite=strict; secure; max-age="+when;
		    document.cookie = "privacy_signer_langs="+navigator.languages.toString()+"; Path=/; secure; SameSite=strict; max-age="+when;
		    $('#modal-privacy-agreement').modal('hide')
		}
		if (!document.cookie.split(';').filter((item) => item.trim().startsWith('privacy_signer_time=')).length) {
			$('#modal-privacy-agreement').modal('show')
		}
	}
}