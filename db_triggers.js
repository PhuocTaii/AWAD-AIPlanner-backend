exports = async function(to, task) {
  const sendgridApiKey = "SendGrid-API key" // API key của SendGrid;
  const sendgridUrl = "https://api.sendgrid.com/v3/mail/send";

  const emailPayload = {
    personalizations: [{
      to: [{ email: to }],
      subject: `Task "${task.name}" đã hết hạn`
    }],
    from: { email: "hcmusaisp@gmail.com", name: "AI Planner" },
    content: [{
      type: "text/plain",
      value: `
        Xin chào,

        Task "${task.name}" đã hết hạn vào lúc ${task.estimated_end_time}.
        Vui lòng kiểm tra và xử lý.

        Trân trọng,
        AI Planner
      `
    }]
  };

  try {
    const response = await context.http.post({
      url: sendgridUrl,
      headers: {
        "Authorization": [`Bearer ${sendgridApiKey}`],
        "Content-Type": ["application/json"]          
      },
      body: JSON.stringify(emailPayload)
    });

    if (response.statusCode === 202) {
      return `Email gửi thành công tới ${to}`;
    } else {
      throw new Error(`Gửi email thất bại: ${response.statusCode}`);
    }
  } catch (err) {
    console.log(`Lỗi khi gửi email tới ${to}:`, err.message);
    throw err;
  }
};

exports = async function () {
  // A Scheduled Trigger will always call a function without arguments.
  // Documentation on Triggers: https://www.mongodb.com/docs/atlas/atlas-ui/triggers

  // Functions run by Triggers are run as System users and have full access to Services, Functions, and MongoDB Data.

  // Get the MongoDB service you want to use (see "Linked Data Sources" tab)
  const { Int32 } = BSON; // Lấy Int32 từ BSON
  const serviceName = "mongodb-atlas";
  const databaseName = "AI_Planner";
  const taskcollectionName = "tasks";
  const userCollectionName = "users";
  const taskCollection = context.services
    .get(serviceName)
    .db(databaseName)
    .collection(taskcollectionName);
  const userCollection = context.services
    .get(serviceName)
    .db(databaseName)
    .collection(userCollectionName);

  // Sử dụng context để lấy thư viện BSON
  // const BSON = context.services.get("mongodb-atlas").BSON;

  try {
    const now = new Date();

    // Tìm các tasks đã hết hạn
    const expiredTasks = await taskCollection
      .find({
        status: { $nin: [Int32(3), Int32(2)] }, // Task chưa được đánh dấu 'expired' và 'completed'
        estimated_end_time: { $lt: now },
      })
      .toArray();

    if (expiredTasks.length === 0) {
      console.log("Không có tasks hết hạn.");
      return;
    }

    // Gửi email cho từng task đã hết hạn
    for (const task of expiredTasks) {
      const creatorId = task.user;
      const creator = await userCollection.findOne({ _id: creatorId});
      const creatorEmail = creator.email; // Email người tạo task
      if (!creatorEmail) {
        console.log(`Task ID ${task._id} không có email người tạo.`);
        continue;
      }

      // Gửi email qua hàm sendEmail
      const emailResult = await context.functions.execute(
        "sendEmail",
        creatorEmail,
        task
      );
      console.log(`Đã gửi email tới ${creatorEmail}:`, emailResult);

      // Cập nhật trạng thái task
      await taskCollection.updateOne(
        { _id: task._id },
        { $set: { status: Int32(3) } } // Đánh dấu task là expired
      );
    }
  } catch (err) {
    console.log("Lỗi khi thực hiện cập nhật:", err.message);
  }
};
