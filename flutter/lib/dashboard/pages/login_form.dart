import 'package:dio/dio.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:sdms/client.dart';
import 'package:sdms/const.dart';
import 'package:sdms/dashboard.dart';
import 'package:sdms/dashboard/pages/lecturer_form.dart';
import 'package:sdms/models/lecturer.dart';

class LoginFormPage extends StatefulWidget {
  const LoginFormPage({Key? key, this.next}) : super(key: key);

  String get title => "Staff Login";

  final VoidCallback? next;

  @override
  _LoginFormPageState createState() => _LoginFormPageState();
}

class _LoginFormPageState extends State<LoginFormPage> {
  late final TextEditingController _usernameController;
  late final TextEditingController _passwordController;

  @override
  void initState() {
    _usernameController = TextEditingController();
    _passwordController = TextEditingController();

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.title)),
      body: Container(
        padding: const EdgeInsets.all(kDefaultPadding),
        child: ListView(
          children: [
            Center(
              child: Image.asset(
                "assets/logo.png",
                fit: BoxFit.cover,
                width: 100,
                height: 100,
              ),
            ),
            SizedBox(height: 10),
            Center(child: Text("Scientific Department Management System")),
            SizedBox(height: 10),
            TextField(
              controller: _usernameController,
              decoration: InputDecoration(labelText: "Username"),
            ),
            SizedBox(height: 10),
            TextField(
              controller: _passwordController,
              decoration: InputDecoration(labelText: "Password"),
              obscureText: true,
            ),
            SizedBox(height: 20),
            ElevatedButton.icon(
              onPressed: () async {
                final username = _usernameController.text;
                final password = _passwordController.text;

                if (username.isEmpty || password.isEmpty) {
                  return;
                }

                if (username != "admin" && password != "admin") {
                  try {
                    final res = await client
                        .get("/lecturers/username/$username/$password");
                    final lec = Lecturer.fromJson(res.data);
                    Navigator.pushReplacement(
                        context,
                        MaterialPageRoute(
                            builder: (ctx) => LecturerFormPage(
                                  lecturer: lec,
                                  next: () {},
                                )));
                  } catch (err) {}
                } else {
                  Navigator.pushReplacement(
                      context, MaterialPageRoute(builder: (ctx) => Root()));
                }
              },
              icon: Icon(Icons.vpn_key),
              label: Text("Login to Dashboard"),
            ),
          ],
        ),
      ),
    );
  }
}
