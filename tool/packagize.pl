#!/usr/bin/env perl
use 5.014;
use warnings;
use utf8;
use autodie;
use File::Spec;
use File::Basename 'basename';
use File::Copy 'move';

my $plugin_dir = $ARGV[0];

my $package_name = $plugin_dir =~ s/^check-//r;
   $package_name = "check$package_name";
   $package_name =~ s/-//g;

my $libdir = File::Spec->catfile($plugin_dir, 'lib');
if (-e $libdir) {
    warn "$libdir already exists\n";
    exit;
}
mkdir $libdir;

for my $file (glob "$plugin_dir/*") {
    if ($file =~ /\.md$/ || $file eq 'lib') {
        next;
    }
    if ($file =~ /\.go$/) {
        my $content = slurp_utf8($file);
        $content =~ s/^package main$/package $package_name/ms;
        $content =~ s!func main\(\) \{!// Do the plugin\nfunc Do() {!ms;
        $content =~ s!^// main\n!!ms;
        spew_utf8($file, $content);
    }
    my $base = basename $file;
    my $dst = "$plugin_dir/lib/$base";
    move $file, $dst;
}
my $main = qq[package main

import "github.com/mackerelio/go-check-plugins/$plugin_dir/lib"

func main() {
\t$package_name.Do()
}
];
spew_utf8("$plugin_dir/main.go", $main);

sub slurp_utf8 {
    my $file = shift;
    return do {
        local $/;
        open my $fh, '<:encoding(UTF-8)', $file;
        <$fh>
    };
}

sub spew_utf8 {
    my ($file, $content) = @_;
    open my $fh, '>:encoding(UTF-8)', $file;
    print $fh $content;
}
