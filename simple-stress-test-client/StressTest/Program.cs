using System;
using System.IO;
using System.Net.NetworkInformation;

class Program
{
    static void Main()
    {
        int numPings;
        string numPingsEnv = Environment.GetEnvironmentVariable("num_pings");

        if (!int.TryParse(numPingsEnv, out numPings) || numPings <= 0)
        {
            Console.WriteLine("Invalid value for num_pings environment variable. Using default value: 10");
            numPings = 10;
        }

        int pingInterval;
        string pingIntervalEnv = Environment.GetEnvironmentVariable("ping_interval");

        if (!int.TryParse(pingIntervalEnv, out pingInterval) || pingInterval <= 0)
        {
            Console.WriteLine("Invalid value for ping_interval environment variable. Using default value: 1000ms");
            pingInterval = 1000;
        }

        string logFile = "/tmp/stress-test-log";
        string host = Environment.GetEnvironmentVariable("host"); // Get the host from an environment variable
        if (string.IsNullOrEmpty(host))
        {
            Console.WriteLine("Host not specified. Using default host: www.google.com");
            host = "www.google.com";
        }

        Ping pingSender = new Ping();

        try
        {
            using (StreamWriter sw = new StreamWriter(logFile))
            {
                Console.SetOut(sw);

                for (int i = 0; i < numPings; i++)
                {
                    DateTime currentTime = DateTime.Now; // Get the current system time
                    PingReply reply = pingSender.Send(host);

                    if (reply.Status == IPStatus.Success)
                    {
                        Console.WriteLine($"{currentTime}: Ping {i + 1} to {host} was successful.");
                        Console.WriteLine($"Roundtrip time: {reply.RoundtripTime}ms");
                    }
                    else
                    {
                        Console.WriteLine($"{currentTime}: Ping {i + 1} to {host} failed. Status: {reply.Status}");
                    }

                    if (i < numPings - 1)
                    {
                        System.Threading.Thread.Sleep(pingInterval);
                    }
                }
            }
        }
        catch (PingException ex)
        {
            using (StreamWriter sw = new StreamWriter(logFile))
            {
                Console.SetOut(sw);
                Console.WriteLine($"An exception occurred while pinging {host}: {ex.Message}");
            }
        }
    }
}
