use strict;
use warnings;
use utf8;
use Test::More;

use JSON::PP qw/decode_json/;

my $config = decode_json do {
    local $/;
    open my $fh, '<', 'packaging/config.json' or die $!;
    <$fh>
};
ok $config->{description};

my $plugins_to_be_packaged = $config->{plugins};
isa_ok $plugins_to_be_packaged, 'ARRAY';

my %plugins = map {s/^check-//; ($_ => 1)} <check-*>;
for my $plug (@$plugins_to_be_packaged) {
    ok $plugins{$plug}, "$plug ok";
}

done_testing;
