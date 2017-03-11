package MCUtil;

use 5.014;
use warnings;
use utf8;

use IPC::Cmd qw/run/;
use Carp qw/croak/;

use parent 'Exporter';

our @EXPORT = qw/
    command_with_exit_code command git hub
    http_get
    debugf infof warnf errorf
    replace slurp spew
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
