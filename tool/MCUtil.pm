package MCUtil;

use 5.014;
use warnings;
use utf8;

use IPC::Cmd qw/run/;
use Carp qw/croak/;
use ExtUtils::MakeMaker qw/prompt/;
use JSON::PP qw/decode_json/;
use version;

use parent 'Exporter';

our @EXPORT = qw/
    command_with_exit_code command git hub
    http_get
    debugf infof warnf errorf
    replace slurp spew
    parse_version decide_next_version suggest_next_version is_valid_version
    last_release merged_prs build_pull_request_body
    scope_guard/;

sub DEBUG() { $ENV{MC_RELENG_DEBUG} }

sub command_with_exit_code {
    say('+ '. join ' ', @_) if DEBUG;
    my $ret = system(@_);
}

sub command {say('+ '. join ' ', @_) if DEBUG; !system(@_) or croak $!}
sub _git {
     state $com = do {
        chomp(my $c = `which git`);
        die "git command is required\n" unless $c;
        $c;
    };
}
sub git {
    unshift  @_, _git; goto \&command
}

sub _hub {
    state $com = do {
        chomp(my $c = `which hub`);
        die "hub command is required\n" unless $c;
        $c;
    };
}
sub hub {
    unshift @_, _hub; goto \&command;
}

sub http_get {
    my $url = shift;
    my ($ok, $err, undef, $stdout) = run(command => [qw{curl -sf}, $url]);
    return {
        success => $ok,
        content => join('', @$stdout),
    };
}

# logger. steal from minilla
use Term::ANSIColor qw(colored);
use constant { LOG_DEBUG => 1, LOG_INFO => 2, LOG_WARN => 3, LOG_ERROR => 4 };

my $Colors = {
    LOG_DEBUG,   => 'green',
    LOG_WARN,    => 'yellow',
    LOG_INFO,    => 'cyan',
    LOG_ERROR,   => 'red',
};

sub _printf {
    my $type = pop;
    return if $type == LOG_DEBUG && !DEBUG;
    my ($temp, @args) = @_;
    my $msg = sprintf($temp, map { defined($_) ? $_ : '-' } @args);
    $msg = colored $msg, $Colors->{$type} if defined $type;
    my $fh = $type && $type >= LOG_WARN ? *STDERR : *STDOUT;
    print $fh $msg;
}

sub infof  {_printf(@_, LOG_INFO)}
sub warnf  {_printf(@_, LOG_WARN)}
sub debugf {_printf(@_, LOG_DEBUG)}
sub errorf {
    my(@msg) = @_;
    _printf(@msg, LOG_ERROR);

    my $fmt = shift @msg;
    die sprintf($fmt, @msg);
}

# file utils
sub slurp {
    my $file = shift;
    local $/;
    open my $fh, '<:encoding(UTF-8)', $file or die $!;
    <$fh>
}
sub spew {
    my ($file, $data) = @_;
    open my $fh, '>:encoding(UTF-8)', $file or die $!;
    $data .= "\n" if $data !~ /\n\z/ms;
    print $fh $data;
}
sub replace {
    my ($file, $code) = @_;
    my $content = $code->(slurp($file));
    spew($file, $content);
}

## version utils
sub parse_version {
    my $ver = shift;
    my ($major, $minor, $patch) = $ver =~ /^([0-9]+)\.([0-9]+)\.([0-9]+)$/;
    ($major, $minor, $patch)
}

sub suggest_next_version {
    my $ver = shift;
    my ($major, $minor, $patch) = parse_version($ver);
    join '.', $major, ++$minor, 0;
}

sub is_valid_version {
    my $ver = shift;
    my ($major) = parse_version($ver);
    defined $major;
}

sub decide_next_version {
    my $current_version = shift;
    my $next_version = suggest_next_version($current_version);
    $next_version = prompt("next version", $next_version);

    if (!is_valid_version($next_version)) {
        die qq{"$next_version" is invalid version string\n};
    }
    if (version->parse($next_version) < version->parse($current_version)) {
        die qq{"$next_version" is smaller than current version "$current_version"\n};
    }
    $next_version;
}

## git utils
sub last_release {
    my @out = `git tag`;

    my ($tag) =
        sort { version->parse($b) <=> version->parse($a) }
        map {/^v([0-9]+(?:\.[0-9]+){2})$/; $1 || ()}
        map {chomp; $_} @out;
    $tag;
}

sub merged_prs {
    my $current_tag = shift;

    my $data = eval { decode_json scalar `ghch -f v$current_tag` };
    if ($! || $@) {
        die "parse json failed: $@";
    }
    return grep {$_->{title} !~ /\[nitp?\]/i} @{ $data->{pull_requests} };
}

sub build_pull_request_body {
    my ($next_version, @releases) = @_;
    my $body = "Release version $next_version\n\n";
    for my $rel (@releases) {
        $body .= sprintf "- %s #%s\n", $rel->{title}, $rel->{number};
    }
    $body;
}

# scope_guard
package MCUtil::g {
    sub new {
        my ($class, $code) = @_;
        bless $code, $class;
    }
    sub DESTROY {
        my $self = shift;
        $self->();
    }
}
sub scope_guard(&) {
    my $code = shift;
    MCUtil::g->new($code);
}

1
