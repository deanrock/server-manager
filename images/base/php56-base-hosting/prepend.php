<?php
//REDIRECT_REDIRECT_AJAX => REDIRECT_AJAX
$find = 'REDIRECT_REDIRECT_';
foreach($_SERVER as $var => $value) {
    if(strncmp($find, $var, strlen($find)) == 0) {
        $_SERVER[str_replace("REDIRECT_REDIRECT_", "REDIRECT_", $var)] = $value;
    }
}
