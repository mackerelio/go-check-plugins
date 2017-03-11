use strict;
use warnings;
use utf8;
use Test::More;

use FindBin;
use lib $FindBin::Bin;
use MCUtil;

my $version = '0.1.2';
my ($major, $minor, $patch) = parse_version $version;
is $major, 0;
is $minor, 1;
is $patch, 2;
is suggest_next_version($version), '0.2.0';

done_testing;
