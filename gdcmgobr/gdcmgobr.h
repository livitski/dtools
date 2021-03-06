#ifndef GDCMGOBR_H
#define GDCMGOBR_H

#include <string>
#include <gdcmCompositeNetworkFunctions.h>
#include <vector>

bool CEcho (std::string remote, int portno, std::string aetitle, std::string call);
bool CStore (std::string remote, int portno, std::string aetitle, std::string call,std::string files);
std::string CFind(std::string callingaetitle,std::string callaetitle,std::string hostname,int port ,std::string  StudyInstanceUID,std::string PatientName,std::string AccessionNumber,std::string PatienDateOfBirth,std::string StudyDate);
bool CGet(std::string aetitle,std::string call,std::string hostname,int port ,std::string  StudyInstanceUID,	std::string PatientName,std::string AccessionNumber,std::string PatienDateOfBirth,std::string StudyDate,std::string SFolder);
#endif
