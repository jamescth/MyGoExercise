#!/usr/bin/env python
"""
Objective: 
Send automatic notifiacation email via office 365 to applications@omnitier.com

Warning:
This file will not import any 3rd party module in order to reduce the distro dependencies
since all Linux distros have Python as part of the default, this code should be 'plug & play'

Note:
Golang seems to have trouble w/ SMTP over TLS when connecting to Office 365.  Instead of 
wasting time to trouble shoot the Golang issue, using Python as a quick and dirty patch for
now will serve the need.

golang issue
https://github.com/golang/go/issues/9899

"""

import smtplib
import base64

import getopt
import sys
import json

import mimetypes
from email import encoders
from email.utils import COMMASPACE, formatdate

from email.MIMEMultipart import MIMEMultipart
from email.MIMEMessage import MIMEMessage
from email.MIMEImage import MIMEImage
from email.MIMEAudio import MIMEAudio
from email.MIMEBase import MIMEBase
from email.MIMEText import MIMEText

from email.mime.application import MIMEApplication

class Mail(object):
    """
    Class for email content
    """
    def __init__(self):
        self.user = 'Automatic.Notification@omnitier.com'
        self.password = base64.b64decode("WW9rdTQxMzg=")
        self.server = "smtp.office365.com:587"
        self.sendto = None
        self.subject = None
        self.body = None
        self.flist = []

def parsing(argv, mail):
	"""
	parse the input arguments
	"""

	attach_file = None

	if (len(argv) < 2):
		usage()
		sys.exit(1)

        '''
        getopt()
            h: help.  no parameter needed.
            b: attach filename.  use '-b <filename>' or '-b<filename>'
            attach: attach filename (same as -a).  use '--attach=<filename>' or '--attach <filename>'
        '''
	try:
            opts, args = getopt.getopt(argv,"hb:l:t:s:",
					["body=","flist="])

	except getopt.GetoptError:
		usage()
		sys.exit(1)

	for opt, arg in opts:
		if opt == '-h':
		    usage()
		    sys.exit(0)
                elif opt == '-t':
                    mail.sendto = arg.split(',')
                    print mail.sendto
                elif opt == '-s':
                    mail.subject = arg
                elif opt in ("-b", "--body"):
                    mail.body = arg
		elif opt in ("-l", "--flist"):
		    mail.flist = arg.split(',')
                    print mail.flist
	return 


def usage():
	print('   sendmail.py -t <send_to@mail.com> -s subject -a <attach file>')
	print('   py_scapy.py --ifile=<input file> --ofile=<output file>')

def sendmail(mail):
    msg = MIMEMultipart()

    msg['FROM'] = mail.user
    msg['To'] = COMMASPACE.join(mail.sendto)
    msg['Date'] = formatdate(localtime=True)
    msg['Subject'] = mail.subject

    if mail.body is not None:
        msg.attach(MIMEText(file(mail.body).read()))

    if mail.flist is not None:
        print mail.flist
        for fileToSend in mail.flist:
            ctype, encoding = mimetypes.guess_type(fileToSend)
            if ctype is None or encoding is not None:
                ctype = "application/octet-stream"
    
            maintype, subtype = ctype.split("/", 1)

            if maintype == "text":
                fp = open(fileToSend)
                # Note: we should handle calculating the charset
                attachment = MIMEText(fp.read(), _subtype=subtype)
                fp.close()
            elif maintype == "image":
                fp = open(fileToSend, "rb")
                attachment = MIMEImage(fp.read(), _subtype=subtype)
                fp.close()
            elif maintype == "audio":
                fp = open(fileToSend, "rb")
                attachment = MIMEAudio(fp.read(), _subtype=subtype)
                fp.close()
            else:
                fp = open(fileToSend, "rb")
                attachment = MIMEBase(maintype, subtype)
                attachment.set_payload(fp.read())
                fp.close()
                encoders.encode_base64(attachment)

            attachment.add_header("Content-Disposition", "attachment", filename=fileToSend)
            msg.attach(attachment)

    smtpserver = smtplib.SMTP(mail.server)
    #smtpserver.ehlo()
    smtpserver.starttls()
    #smtpserver.ehlo
    smtpserver.login(mail.user, mail.password)

    ## simple version
    # header = 'To:' + sendto + '\n' + 'From: ' + user + '\n' + 'Subject:testing \n'
    # print header
    # msgbody = header + '\n This is a test Email send using Python \n\n'
    # smtpserver.sendmail(user, sendto, msgbody)

    # w/ msg
    smtpserver.sendmail(mail.user, mail.sendto, msg.as_string())
    print msg.as_string()

    print 'done!'
    smtpserver.close()

def main():
    mail = Mail()

    attachedF = parsing(sys.argv[1:], mail)

    # sendto = 'applications@omnitier.com'

    sendmail(mail)

if __name__ == '__main__':
	main()
