module github.com/sufficit/sufficit-quepasa-fork/controllers

require (
	github.com/sufficit/sufficit-quepasa-fork/library v0.0.0-00010101000000-000000000000 // indirect
	github.com/sufficit/sufficit-quepasa-fork/metrics v0.0.0-00010101000000-000000000000 // indirect
    github.com/sufficit/sufficit-quepasa-fork/models v0.0.0-00010101000000-000000000000 // indirect
    github.com/sufficit/sufficit-quepasa-fork/whatsapp v0.0.0-00010101000000-000000000000 // indirect
    github.com/sufficit/sufficit-quepasa-fork/whatsmeow v0.0.0-00010101000000-000000000000 // indirect
)

replace github.com/sufficit/sufficit-quepasa-fork/library => ../library
replace github.com/sufficit/sufficit-quepasa-fork/metrics => ../metrics
replace github.com/sufficit/sufficit-quepasa-fork/models => ../models
replace github.com/sufficit/sufficit-quepasa-fork/whatsapp => ../whatsapp
replace github.com/sufficit/sufficit-quepasa-fork/whatsmeow => ../whatsmeow

go 1.17