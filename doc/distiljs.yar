rule linked_distiljs {
    strings:

        $script_tag = /\<script type=\"text\/javascript\" src=\"\/[0-9a-zA-Z]+\.js\" defer\>/

    condition:
        $script_tag
}

rule distiljs {
	strings:
		$pow = "h(\"DistilPostResponse\")"

	condition:
		$pow

}