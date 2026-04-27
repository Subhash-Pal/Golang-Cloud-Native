$udpClient = New-Object System.Net.Sockets.UdpClient
$udpClient.Connect("localhost", 8081)
$encodedMessage = [System.Text.Encoding]::ASCII.GetBytes("Hello, UDP server! New Client")
$udpClient.Send($encodedMessage, $encodedMessage.Length)
$udpClient.Close()
