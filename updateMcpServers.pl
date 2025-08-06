#!/usr/bin/perl

use strict;

if (scalar @ARGV != 2) {
    die "Usage: $0 <input-gemini-file>  <output-gemini-file> " . scalar(@ARGV)
}

my $inputGeminiExtFile  = shift @ARGV;
my $outputGeminiExtFile = shift @ARGV;

open(IF, $inputGeminiExtFile) || die "Cannot read file: $inputGeminiExtFile";
open(OF, ">$outputGeminiExtFile") || die "Cannot create file: $outputGeminiExtFile";

while (<IF>) {
    if (/mcpServers/) {
        print OF "    \"mcpServers\": { \n";
        print OF "        \"context7\": { \n";
        print OF "             \"httpUrl\": \"https://mcp.context7.com/mcp\" \n";
        print OF "         }, \n";
    } else {
        print OF $_;
    }     
}

close(IF);
close(OF);
